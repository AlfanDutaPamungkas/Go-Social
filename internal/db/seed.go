package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/AlfanDutaPamungkas/Go-Social/internal/store"
)

var usernames = []string{
	"mona", "changli", "zeno", "kairo", "livia", "torin", "selva", "vexon", "dario", "felix",
	"nilo", "soren", "yuna", "tova", "brin", "zelka", "joris", "vanya", "xela", "quin",
	"ravi", "zuma", "elric", "fenna", "galo", "haven", "indra", "juno", "kael", "lexa",
	"miko", "nova", "oriel", "pax", "quora", "rynn", "syra", "thane", "ulric", "vela",
	"wyn", "xion", "yuki", "ziv", "astra", "brio", "casper", "dael", "evren", "fael",
	"gael", "halcy", "ion", "jarek", "kyro", "lyric", "mira", "nyx", "odin", "pyra",
	"quent", "riven", "sable", "tryst", "ursa", "vero", "wynn", "xander", "yara", "zeke",
	"arlo", "blix", "cato", "draven", "elio", "frey", "gide", "haze", "imri", "jax",
	"kian", "laz", "maev", "neo", "ollin", "phel", "quinz", "rory", "silas", "taz",
	"ulyx", "vale", "wolfe", "xylo", "yoran", "zade",
}

var titles = []string{
	"The Daily Muse", "Wander & Wonder", "Echoes of Life", "Beyond the Horizon", "Notes & Thoughts",
	"The Quiet Observer", "Stories Untold", "Moments in Time", "Reflections & Realities", "The Curious Mind",
	"Pathways & Perspectives", "Through My Eyes", "Chasing the Sun", "Bits of Wisdom", "The Honest Voice",
	"Unwritten Chapters", "Fragments of Life", "Between the Lines", "The Open Journal", "Lost & Found",
}

var contents = []string{
	"Every day brings new inspiration that drives us to keep creating.",
	"Life's journey is filled with wonders waiting to be discovered.",
	"Every step we take leaves a mark that echoes through time.",
	"A new world awaits, just beyond the horizon, ready to be explored.",
	"Every thought is a piece of the puzzle that shapes our view of the world.",
	"Sometimes, silence is the best way to understand everything around us.",
	"Everyone has a story worth telling, waiting to be heard.",
	"Time is the most precious thing we have, full of precious moments.",
	"Reflection on life brings us to a deeper understanding of reality.",
	"A curious mind always uncovers new and amazing things.",
	"Life is about choosing paths, and our perspective shapes the direction.",
	"Seeing the world from a personal perspective adds a unique color to every story.",
	"Chasing the sun is a metaphor for seeking endless possibilities in life.",
	"In the stillness of the night, we find our most profound thoughts.",
	"Every challenge is an opportunity for growth and self-discovery.",
	"Simplicity is the key to clarity in a world full of complexity.",
	"The little things often hold the greatest significance in our lives.",
	"Through failure, we find the strength to rise and try again.",
	"The unknown is full of potential, waiting for us to take the first step.",
	"True freedom comes when we let go of all expectations and just be.",
}

var tags = []string{
	"inspiration", "creativity", "mindfulness", "wanderlust", "selfgrowth",
	"motivation", "lifestyle", "reflection", "positivity", "exploration",
	"wellness", "journaling", "adventure", "authenticity", "change",
	"empowerment", "perspective", "mindset", "personaljourney", "success",
}

var comments = []string{
	"Great read!", "Love this perspective!", "So inspiring!", "Well written!", "Totally agree!", 
	"This made my day!", "Thanks for sharing!", "Such a fresh take!", "I needed this today!", "Absolutely true!", 
	"Amazing insights!", "Very thought-provoking!", "Couldn't agree more!", "Well said!", "This resonates deeply!", 
	"Brilliantly written!", "Keep up the great work!", "This is so relatable!", "Insightful and inspiring!", "You nailed it!", 
	"A beautiful reflection!", "This really speaks to me!", "Such a powerful message!", "I love your writing style!", "Deep and meaningful!", 
	"Simply wonderful!", "This got me thinking!", "A must-read!", "Well put together!", "I appreciate this perspective!", 
	"You expressed it perfectly!", "A great reminder!", "Thanks for this wisdom!", "Very eye-opening!", "Pure brilliance!", 
	"Short but impactful!", "Such a deep thought!", "Well articulated!", "This gave me chills!", "Such a refreshing view!", 
	"Beautifully expressed!", "A great takeaway!", "I can relate so much!", "Exactly what I needed!", "This is gold!", 
	"So simple yet profound!", "Absolutely beautiful!", "This hits home!", "Really made me reflect!", "So true and real!", 
	"This blew my mind!", "This made me smile!", "Definitely worth reading!", "Exceptionally written!", "A dose of positivity!", 
	"I learned something new!", "Such an important message!", "Very engaging!", "This touched my heart!", "You captured it so well!", 
	"So beautifully written!", "I love your insights!", "This spoke to my soul!", "Absolutely inspiring!", "Really enjoyed this!", 
	"This is pure wisdom!", "You always inspire me!", "I needed this reminder!", "Such a unique take!", "This makes so much sense!", 
	"Every word resonates!", "This is next-level thinking!", "A fresh perspective!", "Your words are powerful!", "Such clarity in your thoughts!", 
	"I felt this deeply!", "Your writing is magical!", "So beautifully put!", "Absolutely loved this!", "This is perfection!", 
	"This post is a gem!", "It’s like you read my mind!", "This left me speechless!", "I’m bookmarking this!", "What a perspective!", 
	"This is truly eye-opening!", "Words to live by!", "This was a joy to read!", "I’m sharing this with friends!", "So deep yet so simple!", 
	"Really got me thinking!", "This deserves more attention!", "I admire your thoughts!", "This is a masterpiece!", "Simply put, brilliant!", 
	"Truly heartfelt!", "This should be framed!", "Words of wisdom!", "I keep coming back to this!", "This is worth reflecting on!", 
	"Very meaningful!", "Absolutely enlightening!", "I appreciate this insight!", "This article is gold!", "Such a lovely thought!", 
	"Keep writing more!", "This uplifted my mood!",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user: ", err)
			return
		}
	}

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post: ", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment: ", err)
			return
		}
	}	

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)
	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@gmail.com",
			Password: "123123",
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment{
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID: posts[rand.Intn(len(posts))].ID,
			UserID: users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return cms
}
