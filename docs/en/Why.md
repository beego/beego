# Design purposes and ideas
People may ask me why I want to build a new web framework rather than use other good ones. I know there are many excellent web frameworks on the internet and almost all of them are open source, and I have my reasons to do this.

Remember when I was writing the book about how to build web applications with Go, I just wanted to tell people what were my valuable experiences with Go in web development, especially I have been working with PHP and Python for almost ten years. At first, I didn't realize that a small web framework can give great help to web developers when they are learning to build web applications in a new programming language, and it also helps people more by studying its source code. Finally, I decided to write a open source web framework called Beego as supporting materiel for my book.

I used to use CI in PHP and tornado in Python, there are both lightweight, so they has following advantages:

1. Save time for handling general problems, I only need to care about logic part.
2. Learn more about languages by studying their source code, it's not hard to read and understand them because they are both lightweight frameworks.
3. It's quite easy to make secondary development of these frameworks for specific purposes.

Those reasons are my original intention of implementing Beego, and used two chapters in my book to introduce and design this lightweight web framework in GO.

Then I started to design logic execution of Beego. Because Go and Python have somewhat similar, I referenced some ideas from tornado to design Beego. As you can see, there is no different between Beego and tornado in RESTful processing; they both use GET, POST or some other methods to implement RESTful. I took some ideas from [https://github.com/drone/routes](https://github.com/drone/routes) at the beginning of designing routes. It uses regular expression in route rules processing, which is an excellent idea that to make up for the default Mux router function in Go. However, I have to design my own interface in order to implement RESTful and use inherited ideas in Python.

The controller is the most important part of whole MVC model, and Beego uses the interface and ideas I said above for the controller. Although I haven't decided to have to design the model part, everyone is welcome to implement data management by referencing Beedb, my another open source project. I simply adopt Go built-in template engine for the view part, but add more commonly used functions as template functions. This is how a simple web framework looks like, but I'll keep working on form processing, session handling, log recording, configuration, automated operation, etc, to build a simple but complete web framework.

- [Introduction](README.md)
- [Installation](Install.md)