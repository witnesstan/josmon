# josmon
Job site monitoring tool
---
## Table of Contents

- [Project Title](#project-title)
- [Description](#description)
- [Features](#features)
- [Installation](#installation)
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
- Configure through TOML file.

## Installation
TODO

## Usage
TODO

## Notes
This project is made to help me locate remote jobs. Also serves as my portfolio. This is my second Golang project after my first "Hello World" script about 5 years ago.
