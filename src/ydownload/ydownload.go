package main

import (
"log"
"net/http"
"os/exec"
"strings"
"io"
)

func getDownloadURL (origURL string) (string, string, error) {
	cmd := exec.Command("youtube-dl","-g","--restrict-filenames","--no-warnings", "-e", origURL)
	d, err := cmd.Output()
	if (err != nil) {
		return "", "", err
	}
	o := strings.TrimSpace(string(d))
	parts := strings.Split (o, "\n")
        title := strings.Replace(parts[0], " ", "_", -1)
        title = strings.Replace(title, ",", "_", -1)
	return strings.TrimSpace (title), strings.TrimSpace(parts[1]), nil
}

func downloadAndStream (w http.ResponseWriter, r *http.Request) {
	o := r.URL.Query().Get("u")
	log.Println (o)
	title, d, err := getDownloadURL (o)
	if err != nil {
		log.Println ("Error getting download url")
		log.Println (err)
		http.Error (w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := http.Get(d)
	if err != nil {
		log.Println ("Error downloading content")
		log.Println (err)
		http.Error (w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set ("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set ("Content-Length", resp.Header.Get("Content-Length"))
	w.Header().Set ("Content-Disposition:", "attachment; filename=" + title + ".mp4")
	defer resp.Body.Close()
	io.Copy (w, resp.Body)
}

func main() {
	http.HandleFunc("/", downloadAndStream)
	log.Fatal(http.ListenAndServe(":8090", nil))
}
