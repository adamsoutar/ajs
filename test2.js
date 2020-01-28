function test () {
  return function hello () {
    console.log("Hello, world!")
  }
}
function test2 () {
  return if (true) { 1 } else { 2 }
}

test2()
