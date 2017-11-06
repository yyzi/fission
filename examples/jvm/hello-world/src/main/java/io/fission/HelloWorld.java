package io.fission;

import spark.Request;
import spark.Response;

public class HelloWorld implements Function {

	public static void main(String[] args) throws Exception {
		Function app = new HelloWorld();
		System.out.println(app.handle(null, null));
	}

	public Object handle(Request request, Response response) throws Exception {
		return "Hello World!";
	}
}
