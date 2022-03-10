for /D %%f in (*) do (
	echo %%f
	pushd %%f
	for /D %%s in (*) do (
		echo %%s
		pushd %%s
		go run ../../../../cmd/compiler/compiler.go %%s.vm
		popd
	)
	popd
)