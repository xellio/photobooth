package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

//
// flags
//
var (
	port           = flag.Int("port", 1122, "the port to run on")
	photoPath      = flag.String("photopath", "", "the path for storing the photos")
	thumbPath      = flag.String("thumbpath", "", "the path for storing the thumbnails")
	lastImageLimit = flag.Int("limit", 100, "maximum amount of images to load")
)

//
// Image struct holds a image and its thumbnail
//
type Image struct {
	Image string `json:"image"`
	Thumb string `json:"thumb"`
}

//
// TemplateVariables struct holds information for the template/frontend
//
type TemplateVariables struct {
	Title     string   `json:"title"`
	LastImage *Image   `json:"lastImage"`
	Count     int      `json:"count"`
	Images    []*Image `json:"images"`
}

//
// main function
//
func main() {
	flag.Parse()

	http.HandleFunc("/", index)
	http.HandleFunc("/gallery", gallery)
	http.HandleFunc("/photo", takePhoto)

	pfs := http.FileServer(http.Dir(*photoPath))
	http.Handle(*photoPath, http.StripPrefix(*photoPath, pfs))

	tfs := http.FileServer(http.Dir(*thumbPath))
	http.Handle(*thumbPath, http.StripPrefix(*thumbPath, tfs))

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("DEBUG: Running on localhost:" + strconv.Itoa(*port))

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))

}

//
// function for index (/) route (GET)
//
func index(w http.ResponseWriter, r *http.Request) {
	if err := executeTemplate(w, "./templates/index.html", "Photobooth"); err != nil {
		log.Println("ERROR: Problem parsing/executing the index template", err)
	}
}

//
// function for gallery (/gallery) route (GET)
//
func gallery(w http.ResponseWriter, r *http.Request) {
	if err := executeTemplate(w, "./templates/gallery.html", "Gallery"); err != nil {
		log.Println("ERROR: Problem parsing/executing the gallery template", err)
	}
}

//
// function for photo (/photo) route (POST)
//
func takePhoto(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	filename := now.Format("20060102120405") + ".jpg"
	log.Println("DEBUG: take photo", filename)

	binary, err := exec.LookPath("./photo.sh")
	if err != nil {
		log.Println("ERROR: Problem with photo.sh", err)
		err = json.NewEncoder(w).Encode(nil)
		if err != nil {
			log.Println("ERROR: sending json", err)
		}
		return
	}
	cmd := exec.Command(binary, filename, *photoPath, *thumbPath)
	_, err = cmd.CombinedOutput()
	if err != nil {
		log.Println("ERROR: Problem executing photo.sh", err)
		err = json.NewEncoder(w).Encode(nil)
		if err != nil {
			log.Println("ERROR: sending json", err)
		}
		return
	}

	templateVars, err := buildTemplateVars("Photobooth")
	if err != nil {
		log.Println("ERROR: Problem rebuilding template variables", err)
		err = json.NewEncoder(w).Encode(nil)
		if err != nil {
			log.Println("ERROR: sending json", err)
		}
	}

	err = json.NewEncoder(w).Encode(templateVars)
	if err != nil {
		log.Println("ERROR: sending json", err)
	}

}

//
// parse template vars to template
//
func executeTemplate(w io.Writer, templatePath string, title string) error {
	t, err := parseTemplate(templatePath)
	if err != nil {
		return err
	}
	templateVars, err := buildTemplateVars(title)
	if err != nil {
		return err
	}

	if err := t.Execute(w, templateVars); err != nil {
		return err
	}
	return nil
}

//
// build the template variables for passing to the template
//
func buildTemplateVars(title string) (*TemplateVariables, error) {

	imageCount, err := getImageCount()
	if err != nil {
		log.Print("ERROR: unable to count images", err)
	}

	lastImages, err := getImages(*lastImageLimit)
	if err != nil {
		log.Println(err)
	}

	lastImage, err := getLastImage()
	if err != nil {
		log.Println("ERROR: unable to fetch last image", err)
	}

	templateVars := &TemplateVariables{
		Title:     title,
		LastImage: lastImage,
		Count:     imageCount,
		Images:    lastImages,
	}

	return templateVars, nil
}

//
// parse the template
//
func parseTemplate(templatePath string) (*template.Template, error) {
	return template.ParseFiles(templatePath)
}

//
// get slice of images limited by count parameter
//
func getImages(count int) ([]*Image, error) {
	lastImages := []*Image{}

	files, err := ioutil.ReadDir(*thumbPath)
	if err != nil {
		return lastImages, err
	}

	for _, f := range files {
		if len(lastImages) >= count {
			break
		}

		lastImages = append(lastImages, &Image{
			Image: *photoPath + f.Name(),
			Thumb: *thumbPath + f.Name(),
		})
	}
	return lastImages, nil
}

//
// get the last (newest) image
//
func getLastImage() (*Image, error) {
	images, err := getImages(1)
	if err != nil || len(images) <= 0 {
		return &Image{}, err
	}

	return images[0], nil
}

//
// count existing images
//
func getImageCount() (int, error) {
	files, err := ioutil.ReadDir(*photoPath)
	if err != nil {
		return 0, err
	}
	return len(files) / 2, nil
}
