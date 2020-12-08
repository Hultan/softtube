module github.com/hultan/softtube

go 1.14

// FOR TESTING
//replace github.com/hultan/softteam-tools => ../softteam-tools

require (
	github.com/go-sql-driver/mysql v1.5.0
	// 201027 : We use latest version here instead of v0.4.0,
	// since gdk.PixbufGetType() was accidently removed in v0.4.0
	github.com/gotk3/gotk3 v0.4.1-0.20200919055744-5c37c8051f06
	//github.com/gotk3/gotk3 v0.5.1
	github.com/hultan/softteam v0.1.2
	github.com/hultan/softteam-tools v1.2.4
)
