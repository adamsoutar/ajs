function test1 (func) {
    return function test2 () {
        return func()
    }
}

test1(function testing () {
    console.log("Hello, world!")
})()