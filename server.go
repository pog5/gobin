package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", respond)
	http.ListenAndServe(":80", nil)
}

func respond(w http.ResponseWriter, r *http.Request) {
	// create
	if r.Method == "POST" {
		paste, err := io.ReadAll(r.Body)
		if len(paste) > 1e7 {
			http.Error(w, "File too large (Over 10MB)", http.StatusRequestEntityTooLarge)
			return
		}
		if err != nil {
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}
		randomNumber := fmt.Sprint(rand.Intn(1000) + 1)
		pastePath := fmt.Sprint("pastes/" + randomNumber)
		_ := os.Mkdir("pastes", os.ModePerm)
		os.WriteFile(pastePath, paste, 0644)
		fmt.Fprintf(w, `{"key": "%s"}`, randomNumber)
		return
	}

	// ui
	if r.URL.Path == "/" {
		w.Header().Set("Content-Type", "text/html")
		page := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		    <meta charset="utf-8">
		    <title>HTTP GoBin</title>
		    <style>
		        * {
		            box-sizing: border-box;
		            font-family: Monaco, Consolas, monospace;
		            font-size: 1.0em;
		        }
		        #content {
		            width: 100%;
		            height: 100%;
		            padding: .5em;
		            border: 0;
		            margin: 0;
		            resize: none;
		            color: #e0e2e4;
		            background: #283134;
		        }
		        #content:focus {
		            outline: none;
		        }
		        #content-wrapper {
		            position: absolute;
		            top: 2em;
		            left: 0;
		            right: 0;
		            bottom: 0;
		        }
		        #nav {
		            position: absolute;
		            left: 0;
		            right: 0;
		            top: 0;
		            height: 2em;
		            line-height: 2em;
		            text-align: right;
		            vertical-align: middle;
		            background: #232323;
		            color: #f97705;
		            white-space: nowrap;
		        }
		        #nav .nav-left {
		            float: left;
		        }
		        #nav a {
		            display: inline-block;
		            padding: 0 .5em;
		            text-decoration: none;
		            color: inherit;
		            cursor: pointer;
		        }
		        #nav a:hover {
		            background: #666666;
		        }
		    </style>
		    <script src="https://cdnjs.cloudflare.com/ajax/libs/pako/2.0.3/pako_deflate.min.js" integrity="sha512-1LtY6ivTYdyp7yVk1N3ZW2wHWT+nA36dWhLNmg5FORFhEAMPdDXG30E1KeeZDYO9659MYjQBIqWEL9EnNOBb4w==" crossorigin="anonymous"></script>
		    <script type="text/javascript">
		        function submit() {
		            var content = document.getElementById("content").value;
		            if (!content) {
		                return
		            }
		            var compressed = pako.gzip(content)
		            var xhr = new XMLHttpRequest();
		            xhr.open("POST", ".", true);
		            xhr.setRequestHeader("Content-Type", "text/plain");
		            xhr.setRequestHeader("Content-Encoding", "gzip");
					xhr.send(compressed);
		            xhr.onload = () => {
		                const fileKey = JSON.parse(xhr.responseText).key;
						window.location.href = "/pastes/" + fileKey;
		            }
		        }
		    </script>
		</head>
		<body>
			<div id="nav">
		    	<a class="nav-left" onclick="submit()">[upload]</a>
		    	<a href="https://github.com/pog5/gobin">[about]</a>
			</div>
			<div id="content-wrapper">
		    	<textarea autofocus id="content" spellcheck="false" placeholder="Type or paste your content here, then click 'upload'..."></textarea>
			</div>
		</body>
		</html>
		`
		fmt.Fprint(w, page)
		return
	}

	// read
	url := string(r.URL.Path)
	path := url[1:]
	fi, err := os.Open(string(path))
	if err != nil {
		errmsg := fmt.Errorf("error reading file: %v", err)
		http.Error(w, string(errmsg.Error()), http.StatusInternalServerError)
		return
	}
	fz, _ := gzip.NewReader(fi)
	s, _ := io.ReadAll(fz)
	fi.Close()
	fmt.Fprint(w, string(s))
}
