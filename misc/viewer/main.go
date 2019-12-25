package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
)

type Type string

const (
	Connd = Type("connd")
	Pushd = Type("pushd")
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":8080", "listen address")
	flag.Parse()

	router := gin.Default()

	router.LoadHTMLGlob("templates/*.tmpl")
	router.Static("/static", "static")
	router.Static("/tempfiles", "./tempfiles")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	router.GET("/", index)
	router.GET("/index", index)
	//router.GET("/draw", draw)

	fmt.Printf("start server and listen in %s\n", addr)
	router.Run(addr)
}

type Data struct {
	ConndList string
	PushdList string
	ErrorMsg  string
	ImagePath string
}

func index(c *gin.Context) {
	var data Data
	data.ConndList = c.Query("connd-addr")
	data.PushdList = c.Query("pushd-addr")
	if len(data.ConndList) == 0 || len(data.PushdList) == 0 {
		c.HTML(http.StatusOK, "index.tmpl", data)
		return
	}
	conndList, err := parse(Connd, data.ConndList)
	if err != nil {
		data.ErrorMsg = err.Error()
		c.HTML(http.StatusOK, "index.tmpl", data)
		return
	}
	pushdList, err := parse(Pushd, data.PushdList)
	if err != nil {
		data.ErrorMsg = err.Error()
		c.HTML(http.StatusOK, "index.tmpl", data)
		return
	}
	list := append(conndList, pushdList...)
	grap := draw(list)
	path, err := genPng(grap, "")
	if err != nil {
		data.ErrorMsg = err.Error()
		c.HTML(http.StatusOK, "index.tmpl", data)
		return
	}
	data.ImagePath = path
	c.HTML(http.StatusOK, "index.tmpl", data)
}

type Node struct {
	Name string
	Line []string
}

func parse(t Type, str string) ([]*Node, error) {
	list, err := split(str)
	if err != nil {
		return nil, err
	}
	var nodes []*Node
	for i := range list {
		node, err := getNode(t, list[i])
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func split(str string) ([]string, error) {
	var list []string
	if strings.Contains(str, ",") {
		list = strings.Split(str, ",")
	} else {
		list = strings.Split(str, "\n")
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("parse list failed. %s", str)
	}
	var res []string
	for i := range list {
		l := strings.Trim(list[i], "[ ,\n")
		if len(l) > 0 {
			res = append(res, l)
		}
	}
	return res, nil
}

func getNode(t Type, addr string) (*Node, error) {
	switch t {
	case Connd:
		return getConndNode(addr)
	case Pushd:
		return getPushdNode(addr)
	default:
	}
	return nil, fmt.Errorf("Not found type")
}

func getConndNode(addr string) (*Node, error) {
	body, err := getBody(addr)
	if err != nil {
		return nil, err
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}
	listen, err := js.GetPath("conf", "Grpc", "Listen").String()
	if err != nil {
		return nil, err
	}
	node := &Node{
		Name: listen,
	}
	pushdCluster := js.Get("pushdCluster")
	arr, err := pushdCluster.Array()
	if err != nil {
		return nil, err
	}
	if len(arr) == 0 {
		return nil, fmt.Errorf("not get pushdCluster array")
	}
	for i := 0; i < len(arr); i++ {
		reg := pushdCluster.GetIndex(i)
		nodes := reg.Get("nodes")
		narr, err := nodes.Array()
		if err != nil {
			return nil, err
		}
		for j := 0; j < len(narr); j++ {
			node.Line = append(node.Line, nodes.GetIndex(i).Get("addr").MustString())
		}
	}
	return node, nil
}

func getPushdNode(addr string) (*Node, error) {
	body, err := getBody(addr)
	if err != nil {
		return nil, err
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}
	listen, err := js.GetPath("conf", "Listen").String()
	if err != nil {
		return nil, err
	}
	node := &Node{
		Name: listen,
	}
	nodes := js.GetPath("conndCluster", "nodes")
	arr, err := nodes.Array()
	if err != nil {
		return nil, err
	}
	if len(arr) == 0 {
		return nil, fmt.Errorf("not get connd nodes error")
	}
	for i := 0; i < len(arr); i++ {
		connd := nodes.GetIndex(i)
		node.Line = append(node.Line, connd.Get("addr").MustString())
	}
	return node, nil
}

func getBody(addr string) ([]byte, error) {
	url := fmt.Sprintf("http://%s/debug/vars", addr)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func draw(nodeList []*Node) *gographviz.Graph {
	nodes := make(map[string]bool)
	for _, node := range nodeList {
		nodes[node.Name] = true
		for _, line := range node.Line {
			nodes[line] = true
		}
	}
	graphAst, _ := gographviz.Parse([]byte(`digraph G{}`))
	graph := gographviz.NewGraph()
	gographviz.Analyse(graphAst, graph)

	for k, _ := range nodes {
		graph.AddNode("G", String(k), nil)
	}
	for _, node := range nodeList {
		for _, dest := range node.Line {
			graph.AddEdge(String(node.Name), String(dest), true, nil)
		}
	}
	fmt.Printf("get graph %s\n", graph.String())
	return graph
}

func String(str string) string {
	new := strings.Replace(str, "\"", "\\\"", -1)
	return "\"" + new + "\""
}

func genPng(graph *gographviz.Graph, pngFilePath string) (string, error) {
	tmpfile, err := ioutil.TempFile("./tempfiles", "")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpfile.Name())

	// 输出文件
	if _, err := tmpfile.Write([]byte(graph.String())); err != nil {
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", err
	}
	// 产生图片
	if pngFilePath == "" {
		pngFilePath = tmpfile.Name() + ".png"
	}
	if err = system(fmt.Sprintf("dot %s -T png -o %s", tmpfile.Name(), pngFilePath)); err != nil {
		return "", err
	}
	return pngFilePath, nil
}

//调用系统指令的方法，参数s 就是调用的shell命令
func system(s string) error {
	cmd := exec.Command(`/bin/sh`, `-c`, s) //调用Command函数
	var out bytes.Buffer                    //缓冲字节

	cmd.Stdout = &out //标准输出
	err := cmd.Run()  //运行指令 ，做判断
	if err != nil {
		return err
	}
	return nil
}
