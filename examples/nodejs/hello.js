
module.exports = async function(context) {
    console.log("log message")
    return {
        status: 200,
        body: "Hello, world!\n"
    };
}
