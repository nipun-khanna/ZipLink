package main

import (
	"ZipLink/utils"
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/joho/godotenv"
)

var ctx = context.Background()

func main() {
	godotenv.Load()
	redisCli := utils.NewRedisClient()
	mongoCli, err := utils.NewMongoClient()

	if err != nil {
		fmt.Println("Failed to connect to MongoDB: ", err)
		return
	} else if redisCli == nil {
		fmt.Println("Failed to connect to Redis")
		return
	}

	defer mongoCli.Disconnect(context.Background())

	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(writer, nil)
	})

	http.HandleFunc("/shorten", func(writer http.ResponseWriter, req *http.Request) {
		url := req.FormValue("url")
		shortCode := utils.GetShortCode()
		shortURL := fmt.Sprintf("http://localhost:8080/r/%s", shortCode)

		utils.SetKey(&ctx, redisCli, shortCode, url, 0)
		utils.SaveURLToMongo(ctx, mongoCli, shortCode, url)

		fmt.Fprintf(writer,
			`<div class="flex flex-col items-center bg-green-50 border border-green-200 rounded-xl p-4 shadow mt-4 animate-fade-in">
				<span class="text-green-700 font-semibold text-lg mb-2 flex items-center">
					<svg xmlns='http://www.w3.org/2000/svg' class='h-5 w-5 mr-2' fill='none' viewBox='0 0 24 24' stroke='currentColor'><path stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M9 12l2 2l4-4' /></svg>
					Shortened URL
				</span>
				<a href="/r/%s" class="text-blue-700 font-mono text-base break-all underline hover:text-blue-900 transition">%s</a>
				<button onclick="navigator.clipboard.writeText('%s');this.innerText='Copied!';this.classList.add('bg-green-200');setTimeout(()=>{this.innerText='Copy Link';this.classList.remove('bg-green-200')},1200);" class="mt-3 px-4 py-1 bg-blue-100 text-blue-700 rounded hover:bg-blue-200 transition text-sm">Copy Link</button>
			</div>`, shortCode, shortURL, shortURL)
	})

	http.HandleFunc("/r/{code}", func(writer http.ResponseWriter, req *http.Request) {
		key := req.PathValue("code")
		if key == "" {
			http.Error(writer, "Invalid URL", http.StatusBadRequest)
			return
		}

		longURL, err := utils.GetLongURL(&ctx, redisCli, key)
		if err != nil {
			longURL, err = utils.GetURLFromMongo(ctx, mongoCli, key)
			if err != nil {
				http.Error(writer, "ZipLink not found", http.StatusNotFound)
				fmt.Printf("Error: %v\n", err)
				return
			}
		}
		http.Redirect(writer, req, longURL, http.StatusPermanentRedirect)
	})

	fmt.Println("Server is started - http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
