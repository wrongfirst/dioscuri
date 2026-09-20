{
	// Relative path
	p, err := CanonicalizePath("index.html")
	Tests.EqualCheck("index.html error is nil", nil, err)
	Tests.EqualCheck("index.html normalized", "/index.html", p)
}

{
	// Windows separators and relative components
	p, err := CanonicalizePath(`assets\css\..\css\style.css`)
	Tests.EqualCheck("windows path error is nil", nil, err)
	Tests.EqualCheck("windows path normalized", "/assets/css/style.css", p)
}

{
	// Windows drive letter stripping
	p, err := CanonicalizePath(`C:\var\www\app.js`)
	Tests.EqualCheck("windows drive letter error is nil", nil, err)
	Tests.EqualCheck("windows drive letter stripped", "/var/www/app.js", p)
}

{
	// Leading dot-slash and redundant slashes
	p, err := CanonicalizePath("./favicon.ico")
	Tests.EqualCheck("dot-slash error is nil", nil, err)
	Tests.EqualCheck("dot-slash normalized", "/favicon.ico", p)

	pMulti, errMulti := CanonicalizePath("///assets//img///logo.png//")
	Tests.EqualCheck("multiple slashes error is nil", nil, errMulti)
	Tests.EqualCheck("multiple slashes collapsed and trailing stripped", "/assets/img/logo.png", pMulti)
}

{
	// Empty path and root path
	pEmpty, errEmpty := CanonicalizePath("")
	Tests.EqualCheck("empty path error is nil", nil, errEmpty)
	Tests.EqualCheck("empty path normalizes to root", "/", pEmpty)

	pRoot, errRoot := CanonicalizePath("/")
	Tests.EqualCheck("root path error is nil", nil, errRoot)
	Tests.EqualCheck("root path normalizes to root", "/", pRoot)
}

{
	// Traversal outside root (Security checks)
	_, err := CanonicalizePath("../secret.env")
	Tests.EqualCheck("relative escape attempt returns error", ErrPathEscapesRoot, err)

	_, errNested := CanonicalizePath("assets/../../secret.env")
	Tests.EqualCheck("nested escape attempt returns error", ErrPathEscapesRoot, errNested)

	_, errRootEscape := CanonicalizePath("/../etc/passwd")
	Tests.EqualCheck("root escape returns error", ErrPathEscapesRoot, errRootEscape)

	_, errDeepEscape := CanonicalizePath("/../../../../etc/shadow")
	Tests.EqualCheck("deep escape returns error", ErrPathEscapesRoot, errDeepEscape)
}
