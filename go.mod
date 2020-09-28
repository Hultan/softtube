module github.com/hultan/softtube

go 1.15

// FOR TESTING
//replace github.com/hultan/softteam-tools => ../softteam-tools

require (
	github.com/go-sql-driver/mysql v1.5.0
	// 201027 : We use latest version here, since gdk.PixbufGetType() was accidently removed in v0.4.0
	github.com/gotk3/gotk3 v0.4.1-0.20200919055744-5c37c8051f06
	github.com/hultan/softteam-tools v1.2.2
)
