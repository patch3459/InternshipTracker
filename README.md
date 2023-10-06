## InternshipTracker - Indexing internships passively

This project utilizes Go's concurrency models for very fast internship scraping and updating. This project would work greatly with cron jobs for "passive" scraping, allowing you to spend your time applying to internships instead of finding them. 

### How to run/build. 

1. Ensure you have the Go programming language installed on your computer and added to your path.

2. Clone the repository \
```git clone https://github.com/patch3459/InternshipTracker.git```

3. Cd into the directory \
```cd InternshipTracker```

4. Install dependencies \
``` go get . ```

5. You can run the program by typing \
```go run .```

6. Alternatively, if you'd like to build the program type  
```go build .```

### The Config File -- Config.Json

The InternshipTracker's config file, named config.json, holds information about the spreadsheets that it will read and write from and keywords it will look for. 

### Values in the config file : 

company_list_csv_path : a path leading to a spreadsheet that the tracker will read company job board links from. Set by default to "./JobLinks.csv".

job_list_csv_path : a path leading to a spreadsheet that the tracker will write prospective jobs to. Set by default to "./JobList.csv".

keywords : a list of keywords that the internship tracker will use to find matching jobs. 


### How to add internships for tracking.

The InternshipTracker works by reading from a csv file that holds information about a company's job page and then indexes it for any jobs that match a set of keywords defined in the config.json. 

If you'd like to add job boards to index, go to the csv file defined at company_list_csv in config.json and add a job. 

For id, you can make it whatever you like at the moment. 

For company name, make sure it's name matches the name used in the job board. Case sensitivity counts. 

For listing type, please set it to 1 if the job board is a greenhouse.io link. Set it to 2 if it is a myworkday jobs link. At the moment, we **only support these two**. Do not add other kinds of links at the moment. 

You can find scraped jobs in the csv file defined at  job_list_csv_path in config.json. 

If you are getting started, I have left some sample websites for you to scrape from. I highly recommend these companies. 

### Plans for the future

This project is still being developed. I hope to add more support for more jobs, hashing so that a user can track individual jobs other than websites, a tutorial for using cron jobs, and so much more!

This project is using a Creative Commons Attribution-NonCommercial-ShareAlike 3.0 Unported license. 




