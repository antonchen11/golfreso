// golfreso v1 Saved on 20201214.

package main

import (
	"flag"
	"fmt"
	"github.com/go-rod/rod"
	"log"
	"strings"
	"time"
	"os"
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


// gcTranslate is to map golf course selections.
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

func stringInSlice(a string, list []string) (int, bool) {
	for i, b := range list {
		if b == a {
			return i, true
		}
	}
	return -1, false
}


func main() {

	// Initial mapping of golf course names to ID required to form query request.
	m := map[string]string{"Langara": "1", "Fraserview": "2", "McCleery": "3"}

	accPtr := flag.String("email", "default@gmail.com", "email of COV golf reservations")
	passPtr := flag.String("pass", "password123", "password of COV golf reservations")
	gcPtr := flag.String("gc", "L", "Golf course to reserve. - [1] [L]angara, [2] [F]raserview, [3] [M]cCleery")
	playersPtr := flag.String("p", "4", "Number of players. - 2, 3, 4")
	datePtr := flag.String("d", "2021-01-25", "Date to play - YYYY-MM-DD")
	timePtr := flag.String("t", "8AM", "Time to Play - Morning (8:00 - 12:00), 9AM, 10AM, 11AM, 12AM")
	timeSlotPtr := flag.String("ts", "8:00", "Time to Play - Morning (8:00 - 12:00), 9AM, 10AM, 11AM, 12AM")

	// Time slots increments of every 9 minutes. LOL!
	timeSlice := []string{"07:45", "07:54", "08:03","08:12","08:21","08:30","08:39","08:48","08:57","09:06","09:15","09:24","09:33","09:42","09:51","10:00","10:09","10:18","10:27","10:36","10:45","10:54","11:03","11:12","11:21","11:30","11:39","11:48","11:57","12:06","12:15","12:24","12:33","12:42","12:51","13:00","13:09","13:18","13:27","13:36", "13:45", "13:54", "14:03", "14:12", "14:21","14:30","14:39","14:48","14:57","15:06","15:15","15:24","15:33","15:42","15:51","16:00","16:09","16:18","16:27","16:36","16:45","16:54"}

	//timeSliceWeekends := []string{"07:45", "07:54", "08:03","08:12","08:21","08:30","08:39","08:48","08:57","09:06","09:15","09:24","09:33","09:42","09:51","10:00","10:09","10:18","10:27","10:36","10:45","10:54","11:03","11:12","11:21","11:30","11:39","11:48","11:57","12:06","12:15","12:24","12:33","12:42","12:51","13:00","13:09","13:18","13:27","13:36", "13:45", "13:54", "14:03", "14:12", "14:21","14:30","14:39","14:48","14:57","15:06","15:15","15:24","15:33","15:42","15:51","16:00","16:09","16:18","16:27","16:36","16:45","16:54"}
	fmt.Println("Available time slots:", timeSlice)

	flag.Parse()
	// eval if timeSlotPtr exists in timeSlice or not..
	// https://play.golang.org/p/5BEJ2fenQsV

	k, found := stringInSlice(*timeSlotPtr, timeSlice)
	if found {
		fmt.Printf("%s found at element: %d\n", *timeSlotPtr, k)
		fmt.Println(*timeSlotPtr + " exists!")
	} else {
		fmt.Println("The time slot you have selected " + *timeSlotPtr + " does not exists as a bookable time..")
		fmt.Println("Exiting from application.")
		os.Exit(0)
	}



	// name to gc translator
	pcourseid := gcTranslate(*gcPtr, m)
	fmt.Println("golf course selected: " + *gcPtr + " // "+ pcourseid)

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

	// if isConfirmed is true, build course response.
	if isConfirmed {
		email := *accPtr
		password := *passPtr
		rcourseId := pcourseid
		rplayers := *playersPtr
		rdate :=  *datePtr
		rTimeSlot := *timeSlotPtr
		//rtime := timePtr
		fmt.Println("Building response")
		response := fmt.Sprintf("https://secure.west.prophetservices.com/CityofVancouver/Home/nIndex?CourseId=%s&Date=%s&Time=AnyTime&Player=%s&Hole=18", rcourseId, rdate, rplayers)


		// Login to COV Golf reservations
		fmt.Println("Hit login page")
		//page := rod.New().MustConnect().MustPage("https://www.wikipedia.org/")
		page := rod.New().MustConnect().MustPage("https://secure.west.prophetservices.com/CityofVancouver/Account/nLogOn#Hash")

		fmt.Println("Signing In...")
		page.MustElement("#Email").MustInput(email)
		page.MustElement("#Password").MustInput(password)
		page.MustElement("#frmLogOn > div > div:nth-child(6) > div > a").MustClick()

		//page.MustWaitLoad().MustScreenshot("logingolf.png")
		// If login is successful, in the top right click on the user's profile
		page.MustElement("#SignInNavbarLarge > li.dropdown > a").MustClick()
		fmt.Println("Logged in!")
		// If profile exists in dropdown, Navigate to URL with parameters
		page.MustElement("#SignInNavbarLarge > li.dropdown.open > ul > li:nth-child(2) > a")
		fmt.Println("Logged in and profile found.")
		page.MustNavigate(response)


		// Going too fast, wait for loading of all times available.
		wait := page.MustWaitNavigation()
		wait()
		page.MustWaitLoad()
		fmt.Println("Searching for closest time slot...")

		// Search for elements with the name of rTimeSlot
		// Set a 15-second timeout for all chained actions
		// The total time for search and click must be less than 15 seconds.

		// TIME OUT AND PANIC, NO HANDLING OF PANIC
		//page.Timeout(15 * time.Second).MustSearch(rTimeSlot)

		// Actions after CancelTimeout won't be affected by the 15-second timeout
		//check := func(err error) {
		//	var evalErr *rod.ErrEval
		//	if errors.Is(err, context.DeadlineExceeded) { // timeout error
		//		fmt.Println("Context deadline met, can't find timeslot.. timeout err")
		//	} else if errors.As(err, &evalErr) { // eval error
		//		fmt.Println(evalErr.LineNumber)
		//	} else if err != nil {
		//		fmt.Println("can't handle", err)
		//	}
		//}

		// ATTEMPT 1
		fmt.Println("ATTEMPT 1!")
		fmt.Println("Trying time slot: " + rTimeSlot)
		err := rod.Try(func() {
			//fmt.Println(page.MustElement("a").MustHTML()) // use "Must" prefixed functions
			page.Timeout(15 * time.Second).MustSearch(rTimeSlot)
			// handle timeout
		})
		if err != nil {
			fmt.Println("hmmm... Can't seem to find a time slot for", rTimeSlot, "after 15 seconds...")
		} else {
			page.MustSearch(rTimeSlot).MustElementR("span", rTimeSlot).MustClick()
			fmt.Sprintf("Found time slot of {rTimeSlot} at x golf course")

			// Use card on File
			fmt.Println("Using card on file...")
			page.MustElement("#monerisCardOnFile").MustClick()
			fmt.Println("Reserving!")
			page.MustElement("#btnBook").MustClick()
			time.Sleep(10* time.Second)
			fmt.Printf("Your reservation has been completed! Be ready to play on %s %s at %s \n",rdate, rTimeSlot, *gcPtr)
			fmt.Println("exiting application...")
			os.Exit(0)
		}
		//check(err)

		// ATTEMPT 2 =================================================================
		rTimeSlot2 := timeSlice[k+1]
		fmt.Println("REATTEMPTING - #2")
		fmt.Println("Trying next time slot: " + rTimeSlot2)
;		err1 := rod.Try(func() {
			//fmt.Println(page.MustElement("a").MustHTML()) // use "Must" prefixed functions
			page.Timeout(15 * time.Second).MustSearch(rTimeSlot2)
			// handle timeout
		})
		if err1 != nil {
			fmt.Println("hmmm... Can't seem to find a time slot for", rTimeSlot2, "after 15 seconds...")
		} else {
			page.MustSearch(rTimeSlot2).MustElementR("span", rTimeSlot2).MustClick()
			fmt.Sprintf("Found time slot of {rTimeSlot2} at x golf course")

			// Use card on File
			fmt.Println("Using card on file...")
			page.MustElement("#monerisCardOnFile").MustClick()
			fmt.Println("Reserving!")
			page.MustElement("#btnBook").MustClick()
			time.Sleep(10 * time.Second)
			fmt.Printf("Your reservation has been completed! Be ready to play on %s %s at %s \n",rdate, rTimeSlot, *gcPtr)
			fmt.Println("exiting application...")
			os.Exit(0)
		}

		// ATTEMPT 3 =================================================================
		rTimeSlot3 := timeSlice[k+2]
		fmt.Println("REATTEMPTING - #3 LAST ATTEMPT")
		fmt.Println("Trying next time slot: " + rTimeSlot3)
		err2 := rod.Try(func() {
			//fmt.Println(page.MustElement("a").MustHTML()) // use "Must" prefixed functions
			page.Timeout(15 * time.Second).MustSearch(rTimeSlot3)
			// handle timeout
		})
		if err2 != nil {
			fmt.Println("hmmm... Can't seem to find a time slot for", rTimeSlot3, "after 15 seconds...")
			fmt.Println("exiting application...")
			os.Exit(0)
		} else {
			page.MustSearch(rTimeSlot3).MustElementR("span", rTimeSlot3).MustClick()
			fmt.Sprintf("Found time slot of {rTimeSlot3} at x golf course")

			// Use card on File
			fmt.Println("Using card on file...")
			page.MustElement("#monerisCardOnFile").MustClick()
			fmt.Println("Reserving!")
			page.MustElement("#btnBook").MustClick()
			time.Sleep(10 * time.Second)
			fmt.Printf("Your reservation has been completed! Be ready to play on %s %s at %s \n",rdate, rTimeSlot, *gcPtr)
			fmt.Println("exiting application...")
			os.Exit(0)
		}
 		//check(err1)



		// ===============FEATURE OF ASKING USER FOR TIMESLOT TO TRY.
		//fmt.Println("Do you want to to try another time slot? y/n?")
		//tryAgain := askForConfirmation()
		//
		//if tryAgain {
		//	fmt.Println("Enter time slot to try")
		//
		//	var retryTimeSlot string
		//	fmt.Scanln(&retryTimeSlot)
		//	fmt.Println("Trying for", retryTimeSlot)
		//	page.MustNavigate(response)
		//	page.MustWaitLoad()
		//	fmt.Println("render completed.")
		//	fmt.Println("Searching for new time at ", retryTimeSlot)
		//	page.Timeout(15 * time.Second).MustSearch(retryTimeSlot).MustElementR("span", retryTimeSlot).MustClick()
		//	fmt.Println("Time slot", retryTimeSlot, "found! Continuing to checkout...")
		//} else {
		//	fmt.Println("not found,exiting.")
		//}
		// ==============================================


		// If timeslot not found, retry on next increment of +9 minutes and inform user.
		// Compare rTimeSlot with old and new elements in time array.

		// fmt.Sprintf("The timeslot requested at {rTimeSlot} is not available, attempting next time slot at %s)
		// Attempt x3 times ( est. 27 minute differential of desired time)
		//page.MustSearch(rTimeSlot).MustElementR("span", rTimeSlot).MustClick()
		//fmt.Sprintf("Found time slot of {rTimeSlot} at x golf course")

		//======================== DONT NEED THIS CODE FOR NOW. LOOKS UGLY AF BUT FOR REFERENCE 01/26/2020
		// Confirmation on slot click Y/N?
		//page.MustElement("#enterCCInfo").MustClick()
		// Next
		//page.MustElement("#divSelectPlayersAndHole > div.modal-footer > button.btn.btn-primary.MainNavigationColor.MainNavigationFontColor").MustClick()
		// =====================================================

		// Use card on File
		//page.MustElement("#monerisCardOnFile").MustClick()
		// Book

		// ===================================================== DONT NEED. 01/26/2020
		//page.MustElement("btnBook").MustClick()
		// Credit Card Entry logic. Not required assuming that a credit card is associated with the account already.
		//page.MustElement("avs_str_num").MustClick()
		//page.MustElement("avs_str_name").MustInput("joyce street")
		//page.MustElement("avs_zip_code").MustInput("V2A 2B9")
		//page.MustElement("#cardholder").MustInput("Bob Smith")
		//page.MustElement("#pan").MustInput("123456789")
		//page.MustElement("#exp_month").MustInput("01")
		//page.MustElement("#exp_year").MustInput("2025")
		//page.MustElement("#cvd_value").MustInput("123")
		//page.MustElement("#buttonResAddCC").MustClick()
		//===================================================================


		// Order specific time.
	//	page.MustElementR("span", "08:09").MustClick()
	//	fmt.Println("ordered")
	//	//page.MustElement("#bodyContent > div.container-fluid.teeSheet > div:nth-child(4) > div > a > div > div.col-xs-12.col-md-12.col-lg-12.text-center.divBoxText > div.col-xs-6.col-sm-12.small-padding.p-nopadding > p").MustClick()
	//	time.Sleep(time.Hour)
	//} else {
	//	fmt.Println("Reservation info not confirmed")


	}

}
