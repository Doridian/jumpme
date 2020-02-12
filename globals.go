package main

import "regexp"

var IsRoot bool
var HomeDirs []string
var UIDToName map[int]string

var SshLikeCommand = regexp.MustCompilePOSIX("(ssh|scp|sftp|rsync|git remote|git clone)( |$)")
var HostLikePattern = regexp.MustCompilePOSIX("([a-zA-Z0-9]+@)?[a-zA-Z0-9\\-\\_]+\\.[a-zA-Z0-9\\-\\_\\.]+")
