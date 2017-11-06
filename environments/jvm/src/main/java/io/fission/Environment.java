package io.fission;

import static spark.Spark.delete;
import static spark.Spark.get;
import static spark.Spark.port;
import static spark.Spark.post;
import static spark.Spark.put;

import java.io.File;
import java.net.MalformedURLException;
import java.net.URL;
import java.net.URLClassLoader;
import java.util.Set;

import org.reflections.Reflections;
import spark.Request;
import spark.Response;

public class Environment {

	private static final int HTTP_PORT = 8888;

	private static final File DEFAULT_CODE_PATH = new File("/userfunc/user");

	private Function fn;


	public static void main(String[] args) {
		Environment env = new Environment();
		env.run(HTTP_PORT);
	}

	Reflections loadJar(File jar) {
		try {
			// TODO narrow search to reduce load time
			URLClassLoader child = new URLClassLoader(new URL[] { jar.toURI().toURL() });
			return new Reflections(child); // TODO move to constructor?
		}
		catch (MalformedURLException e) {
			throw new RuntimeException(e);
		}
	}

	private Object specialize(Request req, Response res) throws Exception {
		System.out.println("Specializing environment...");
		Reflections reflect = loadJar(DEFAULT_CODE_PATH);

		// Discover implementations
		Set<Class<? extends Function>> impls = reflect.getSubTypesOf(Function.class);

		// Select right implementation
		Class<? extends Function> impl = impls.iterator().next(); // TODO support multiple implementations

		try {
			fn = impl.newInstance();
		}
		catch (ClassCastException e) {
			res.status(500);
			res.body(String.format("Class could not be cast to %s", Function.class.getCanonicalName()));
		}
		catch (InstantiationException | IllegalAccessException e2) {
			res.status(500);
			res.body(String.format("Could not instantiate user class: %s", impl.getCanonicalName()));
		}
		return "";
	}

	private Object invoke(Request req, Response res) throws Exception {
		if (fn == null) {
			res.status(400);
			res.body("Generic container: no requests supported");
		}

		return fn.handle(req, res);
	}

	public void run(int httpPort) {
		port(httpPort);
		post("/specialize", this::specialize);

		put("/", this::invoke);
		get("/", this::invoke);
		post("/", this::invoke);
		delete("/", this::invoke);
	}
}
