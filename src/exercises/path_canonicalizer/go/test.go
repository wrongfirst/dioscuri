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
	// Leading dot-slash
	p, err := CanonicalizePath("./favicon.ico")
	Tests.EqualCheck("dot-slash error is nil", nil, err)
	Tests.EqualCheck("dot-slash normalized", "/favicon.ico", p)
}

{
	// Traversal outside root
	_, err := CanonicalizePath("../secret.env")
	Tests.EqualCheck("escape attempt returns error", ErrPathEscapesRoot, err)

	_, errNested := CanonicalizePath("assets/../../secret.env")
	Tests.EqualCheck("nested escape attempt returns error", ErrPathEscapesRoot, errNested)

	_, errRootEscape := CanonicalizePath("/../etc/passwd")
	Tests.EqualCheck("root escape returns error", ErrPathEscapesRoot, errRootEscape)
}
