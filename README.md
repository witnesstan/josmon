# josmon
Job site monitoring tool
---
## Table of Contents

- [Project Title](#project-title)
- [Description](#description)
- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Notes](#notes)

## Description
A tool that can detect changes to career web pages and alert the user. This is useful for job seekers who are monitoring a company website for new information regarding the role that they are interested in. The tool is run through a cron job to query the sites regularly. Later, I am going to make this run in serverless like AWS Lambda.

## Features
- Detects changes in HTML pages (doesn't work with JS-loaded content).
- Retrieves list of website URLs from an input file.
- Built-in operation mode to test functions using CLI parameters.
- Cache results to file, retrieve later for compare.
- Send email alerts.
- Uses configuration file.

## Installation
- Grab a copy of the source code.

  `git clone https://github.com/witnesstan/josmon.git`
- Build the code.

  `cd src; go build`
- Run `josmon --help`.

## Configuration
1. Populate the `career_pages.cdf` with the websites you want to monitor. This is a comma-delimited file (not CSV).

   The format is:
   
   &lt;URL&gt;,&lt;boundary_start&gt;,&lt;boundary_end&gt;

   *boundary* refers to the block of html code that you want to focus your search on.
   
   So *boundary_start* is the string (could be multiple words to make it unique) that would mark the start of the boundary.
   
   *boundary_end* is the string that marks the end of the boundary.

   You can't use comma in the *boundary* strings.
3. Configure `josmon.conf`.

   Populate values for the following:
   ```
   email_to <to@addr.com>
   email_from <from@addr.com>
   smtp_user <username>
   smtp_pass <password>
   smtp_host <smtp_host>
   smtp_port 465
   find_keyword <search_keyword>
   ```

   If you don't need to email the results, there's no need to configure `email_*` and `smtp_*` values. But you will surely need to configure `find_keyword`.
## Usage
Run `josmon` without any parameters. It will print its progress.

Run `josmon --help` for help.

## Notes
This project is made to help me locate remote jobs. Also serves as my portfolio. This is my second Golang project after my first "Hello World" script about 5 years ago.
