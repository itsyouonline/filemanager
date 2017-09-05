# ItsYou.Online Integration

This is a fork of the original Filemanager Caddy plugin.

This fork allows you to use the filemanager with itsyouonline authentification and management interface.

## What changed

- The integration is made on the `master-iyo-auth` branch, `master` is kept for upstream update.
- Only a small change on the auth login page was made to the front-end, all the others change are made on the backend
- The database is still used, but it's a temporary file now, and a new one is created on each instance
  - The database contains informations like commands trigger etc. and to avoid lot of code change, keeping it is easier
- We fake the login procedure and we check the jwt provided to the page with itsyouonline privacy settings
- You can add a new option `triggercmd` to the Caddyfile which point to a executable to execute on each (save/edit) action
- The triggercmd will be set with:
  - The working directory contains the directory where the file is edited
  - Environment variables:
    - `file` is the target file
    - `USERNAME` is the itsyouonline username
    - `REALNAME` is the itsyouonline user real name
    - `USEREMAIL` is the itsyouonline user email address
    - `TRIGGER` is the trigger name: `before_save`, `after_save`, `before_publish`, `after_publish`

![Preview](https://user-images.githubusercontent.com/5447088/28537288-39be4288-70a2-11e7-8ce9-0813d59f46b7.gif)

# filemanager

[![Build](https://img.shields.io/travis/hacdias/filemanager.svg?style=flat-square)](https://travis-ci.org/hacdias/filemanager)
[![Go Report Card](https://goreportcard.com/badge/github.com/hacdias/filemanager?style=flat-square)](https://goreportcard.com/report/hacdias/filemanager)
[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/hacdias/filemanager)

filemanager provides a file managing interface within a specified directory and it can be used to upload, delete, preview, rename and edit your files. It allows the creation of multiple users and each user can have its own directory. It can be used as a standalone app or as a middleware.

# Table of contents

+ [Getting started](#getting-started)
+ [Features](#features)
  - [Users](#users)
  - [Search](#search)
+ [Contributing](#contributing)
+ [Donate](#donate)

# Getting started

You can find the Getting Started guide on the [documentation](https://henriquedias.com/filemanager/quick-start/).

# Features

Easy login system.

![Login Page](https://user-images.githubusercontent.com/5447088/28432382-975493dc-6d7f-11e7-9190-23f8037159dc.jpg)

Listings of your files, available in two styles: mosaic and list. You can delete, move, rename, upload and create new files, as well as directories. Single files can be downloaded directly, and multiple files as *.zip*, *.tar*, *.tar.gz*, *.tar.bz2* or *.tar.xz*.

![Mosaic Listing](https://user-images.githubusercontent.com/5447088/28432384-9771bb4c-6d7f-11e7-8564-3a9bd6a3ce3a.jpg)

File Manager editor is powered by [Codemirror](https://codemirror.net/) and if you're working with markdown files with metadata, both parts will be separated from each other so you can focus on the content.

![Markdown Editor](https://user-images.githubusercontent.com/5447088/28432383-9756fdac-6d7f-11e7-8e58-fec49470d15f.jpg)

On the settings page, a regular user can set its own custom CSS to personalize the experience and change its password. For admins, they can manage the permissions of each user, set commands which can be executed when certain events are triggered (such as before saving and after saving) and change plugin's settings.

![Settings](https://user-images.githubusercontent.com/5447088/28432385-9776ec66-6d7f-11e7-90a5-891bacd4d02f.jpg)

We also allow the users to search in the directories and execute commands if allowed.

## Users

We support multiple users and each user can have its own scope and custom stylesheet. The administrator is able to choose which permissions should be given to the users, as well as the commands they can execute. Each user also have a set of rules, in which he can be prevented or allowed to access some directories (regular expressions included!).

![Users](https://user-images.githubusercontent.com/5447088/28432386-977f388a-6d7f-11e7-9006-87d16f05f1f8.jpg)

## Search

FileManager allows you to search through your files and it has some options. By default, your search will be something like this:

```
this are keywords
```

If you search for that it will look at every file that contains "this", "are" or "keywords" on their name. If you want to search for an exact term, you should surround your search by double quotes:

```
"this is the name"
```

That will search for any file that contains "this is the name" on its name. It won't search for each separated term this time.

By default, every search will be case sensitive. Although, you can make a case insensitive search by adding `case:insensitive` to the search terms, like this:

```
this are keywords case:insensitive
```

# Contributing

The contributing guidelines can be found [here](https://github.com/hacdias/filemanager/blob/master/CONTRIBUTING.md).

# Donate

Enjoying this project? You can [donate to its creator](https://henriquedias.com/donate/). He will appreciate.
