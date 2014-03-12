/*global module: false */
module.exports = function(grunt) {
	"use strict";
	grunt.initConfig({
		pkg: grunt.file.readJSON('package.json'),
		compass: {
			dist: {
				options: {
					sassDir: 'static/scss',
					cssDir: 'static/css'
				}
			}
		},
		jshint: {
			options: {
				curly: true,
				eqeqeq: true,
				eqnull: true,
				browser: true,
				loopfunc:true,
				globals: {
					jQuery: true
				},
				ignores:['static/js/*.min.js', 'static/js/libs/**'] // We're not going to jshint libs, too much work :\
			},
			dist: {
				src: ['Gruntfile.js', 'static/js/**/*.js']
			}
		},
		uglify: {
			options: {
				banner: '/* <%= pkg.name %> - version <%= pkg.version %> - ' +
						'<%= grunt.template.today("dd-mm-yyyy") %>\n' +
						'<%= pkg.description %>\n ' +
						'<%= grunt.template.today("yyyy") %> <%= pkg.author.name %> ' +
						'- <%= pkg.author.email %> */\n'
			},
			my_target: {
				files: {
					'static/js/main.min.js': ['static/js/main.js']
				}
			}
		},
		watch: {
			front:{
				files: ['static/scss/*', 'static/js/*.js', 'templates/**/*.html', 'Gruntfile.js'],
				tasks: ['jshint', 'uglify']
			},
			back:{
				options:{
					livereload: true
				},
				files: ['**/**/*.go'],
				tasks:['build-server']
			}
		},
		connect: {
			test: {
				options: {
					port: 9001,
					keepalive: true
				}
			}
		},
		'build-server':{
			dev:{
				root:'$GOPATH/src/github.com/curt-labs/GoAdmin'
			}
		},
		concurrent:{
			options:{
				logConcurrentOutput: true
			},
			prod:{
				tasks:['watch:back','watch:front']
			}
		}
	});

	grunt.loadNpmTasks('grunt-contrib-compass');
	grunt.loadNpmTasks('grunt-contrib-uglify');
	grunt.loadNpmTasks('grunt-contrib-jshint');
	grunt.loadNpmTasks('grunt-contrib-watch');
	grunt.loadNpmTasks('grunt-contrib-jasmine');
	grunt.loadNpmTasks('grunt-concurrent');

	grunt.registerTask('default', ['jshint', 'build-server:dev', 'concurrent:prod']);
	grunt.registerTask('dist', ['jshint','uglify','compass']);

	grunt.registerMultiTask('build-server', "Compile all the things.", function() {
		var done = this.async(),
				path = require('path'),
				dir = grunt.config(['build-server', this.target, 'root']),
				opts = grunt.config(['build-server', this.target, 'opts']);

		return grunt.util.spawn({
			cmd: "go",
			args: ["run", "watcher.go"],
			opts:{stdio:'inherit'}
		}, function (error, result, code) {
			if (error) {
				grunt.log.error("Failed to build server: " + error + " (stdout='" + result.stdout + "', stderr='" + result.stderr + "')");
				done(false);
			} else {
				grunt.log.writeln("Now you can run " + path.resolve(dir, "../", "index") + " and go to http://0.0.0.0:8087");
				done();
			}
		});
	});

};