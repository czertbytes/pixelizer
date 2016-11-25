package handler

import (
	"html/template"
	"net/http"
)

const indexPage = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Pixelizer</title>
        <style>
            .res-w { width: 800px; float: right; }
            img { display: block; width: 100% !important; height: auto !important; }
        </style>
	</head>
	<body>
        <h1>Pixelizer</h1>
        <form id="pix-form">
            <input type="file" id="source" name="file" accept=".jpg,.png">
            <input type="range" id="bs" min=1 max=128 value=16>
            <button name="button" type="submit" id="run">Pixelize</button>
        </form>
        <div class="res-w">
            <img width="500" height="500" id="res" />
        </div>
        <p>Request counter: {{.Counter}}</p>

        <script type="text/javascript">
            var form = document.getElementById('pix-form');
            var runBtn = document.getElementById('run');
            var bs = document.getElementById('bs');

            //bs.oninput = function(e) {
            form.onsubmit = function(e) {
                e.preventDefault();
                runBtn.disabled = true;

                var formData = new FormData();
                formData.append('file', document.getElementById('source').files[0]);

                var xhr = new XMLHttpRequest();
                xhr.responseType = 'blob';
                xhr.onload = function() {
                    if (this.status === 200) {
                        document.getElementById('res').src = window.URL.createObjectURL(this.response);
                    }
                    runBtn.disabled = false;
                };
                xhr.open('POST', '/pixelize?block-size=' + bs.value, true);
                xhr.send(formData);
            }
        </script>
	</body>
</html>`

var indexTpl *template.Template

func init() {
	var err error

	indexTpl, err = template.New("index").Parse(indexPage)
	if err != nil {
		panic(err)
	}
}

type data struct {
	Counter int
}

type Index struct {
	counter int
}

func (h *Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.counter++

		indexTpl.Execute(w, data{h.counter})
		return
	}
}
