// golfreso v1 Saved on 20201214.

package main

import (
	"context"
	"errors"
	"github.com/go-rod/rod"
	"log"
	"strings"
	"time"
	"flag"
	"fmt"
)

func askForConfirmation() bool {
	var response string

	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		fmt.Println("I'm sorry but I didn't get what you meant, please type (y)es or (n)o and then press enter:")
		return askForConfirmation()
	}
}


// IF FIRST TIME IS NOT AVAILABLE THEN
// +/- "3 HOURS DIFF"
// Earliest time you want to play at and see earliest available.
// Idea by Grace.





func gcTranslate(gcname string, m map[string]string) string {
	switch gcname {
	case "Langara","L", "1":
		return m["Langara"]
	case "Fraserview", "F", "2":
		return m["Fraserview"]
	case "McCleery"	, "M", "3":
		return m["McCleery"]
	case "All":
		return "1,2,3"
	}
	return "no courses options recognized"
}

func main() {

	m := map[string]string{"Langara": "1", "Fraserview": "2", "McCleery": "3"}

	accPtr := flag.String("email", "default@gmail.com", "email of COV golf reservations")
	passPtr := flag.String("pass", "password123", "password of COV golf reservations")
	gcPtr := flag.String("gc", "L", "Golf course to reserve. - [1] [L]angara, [2] [F]raserview, [3] [M]cCleery")
	playersPtr := flag.String("p", "4", "Number of players. - 2, 3, 4")
	datePtr := flag.String("d", "2020-12-15", "Date to play - YYYY-MM-DD")
	timePtr := flag.String("t", "8AM", "Time to Play - Morning(8:00 - 12:00), 9AM, 10AM, 11AM, 12AM")
	timeSlotPtr := flag.String("ts", "8:00", "Time to Play - Morning(8:00 - 12:00), 9AM, 10AM, 11AM, 12AM")

	flag.Parse()
	// name to gc translator
	pcourseid := gcTranslate(*gcPtr, m)
	fmt.Println("golf course selected:" + pcourseid)

	fmt.Println("=== VanGOLF reservation information provided ===")
	fmt.Println("Email:", *accPtr)
	fmt.Println("Password:", *passPtr)
	fmt.Println("Golf Course:", *gcPtr)
	fmt.Println("# of Players:", *playersPtr)
	fmt.Println("Date:", *datePtr)
	fmt.Println("Time:", *timePtr)
	fmt.Println("Time Slot:", *timeSlotPtr)
	fmt.Println("=============================================\n")

	fmt.Println("Does the following reservation information look correct? y/n?")
	isConfirmed := askForConfirmation()

	if isConfirmed {
		email := *accPtr
		password := *passPtr
		rcourseId := pcourseid
		rplayers := *playersPtr
		rdate :=  *datePtr
		rTimeSlot := *timeSlotPtr
		//rtime := timePtr
		response := fmt.Sprintf("https://secure.west.prophetservices.com/CityofVancouver/Home/nIndex?CourseId=%s&Date=%s&Time=Morning&Player=%s&Hole=18", rcourseId, rdate, rplayers)


		// Login to COV Golf reservations
		//page := rod.New().MustConnect().MustPage("https://www.wikipedia.org/")
		page := rod.New().MustConnect().MustPage("https://secure.west.prophetservices.com/CityofVancouver/Account/nLogOn#Hash")
		page.MustElement("#Email").MustInput(email)
		page.MustElement("#Password").MustInput(password)
		page.MustElement("#frmLogOn > div > div:nth-child(6) > div > a").MustClick()

		//page.MustWaitLoad().MustScreenshot("logingolf.png")

		// If login is successful, in the top right click on the user's profile
		page.MustElement("#SignInNavbarLarge > li.dropdown > a").MustClick()
		// If profile exists in dropdown, Navigate to URL with parameters
		page.MustElement("#SignInNavbarLarge > li.dropdown.open > ul > li:nth-child(2) > a")
		fmt.Println("Logged in and profile found.")
		page.MustNavigate(response)

		// Going too fast, wait for load.
		wait := page.MustWaitNavigation()
		wait()
		page.MustWaitLoad()
		fmt.Println("Searching for time slot...")
		// Set a 15-second timeout for all chained actions
		// The total time for search and click must be less than 15 seconds.
		//page.Timeout(15 * time.Second).MustSearch(rTimeSlot)
		//page.MustSearch(rTimeSlot).MustElementR("span", rTimeSlot).MustClick()
		// Actions after CancelTimeout won't be affected by the 15-second timeout

		check := func(err error) {
			var evalErr *rod.ErrEval
			if errors.Is(err, context.DeadlineExceeded) { // timeout error
				fmt.Println("timeout err")
			} else if errors.As(err, &evalErr) { // eval error
				fmt.Println(evalErr.LineNumber)
			} else if err != nil {
				fmt.Println("can't handle", err)
			}
		}

		err := rod.Try(func() {
			//fmt.Println(page.MustElement("a").MustHTML()) // use "Must" prefixed functions
			page.Timeout(15 * time.Second).MustSearch(rTimeSlot)
		})
		check(err)

		fmt.Println("hmmm... Can't seem to find a time slot for", rTimeSlot, "after 10 seconds...")
		fmt.Println("Do you want to to try another time slot? y/n?")
		tryAgain := askForConfirmation()

		if tryAgain {
			fmt.Println("Enter time slot to try")

			var retryTimeSlot string
			fmt.Scanln(&retryTimeSlot)
			fmt.Println("Trying for", retryTimeSlot)
			page.MustNavigate(response)
			page.MustWaitLoad()
			fmt.Println("render completed.")
			fmt.Println("Searching for new time at ", retryTimeSlot)
			page.MustSearch(retryTimeSlot).MustElementR("span", retryTimeSlot).MustClick()
			fmt.Println("Time slot", retryTimeSlot, "found! Continuing to checkout...")
		} else {
			fmt.Println("not found,exiting.")
		}



		// Order specific time.
		//page.MustElementR("span", "08:10").MustClick()
		fmt.Println("ordered")
		//page.MustElement("#bodyContent > div.container-fluid.teeSheet > div:nth-child(4) > div > a > div > div.col-xs-12.col-md-12.col-lg-12.text-center.divBoxText > div.col-xs-6.col-sm-12.small-padding.p-nopadding > p").MustClick()
		time.Sleep(time.Hour)
	} else {
		fmt.Println("Reservation info not confirmed")
	}

}
