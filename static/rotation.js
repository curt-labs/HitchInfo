rotateInterval = ""
rotationDuration = 10000;
$(document).ready(function() {
	//Load the slideshow
	sizeThumbs();
	$('#rotation div.panel').hoverIntent(function(){pauseRotation()},function(){playRotation()})
	$('#forward').click(function(){if($('div#rotation div.panel').length > 1) rotate(1)})
	$('#back').click(function(){if($('div#rotation div.panel').length > 1) rotate(-1)})
	//Set the opacity of all images to 0
	//This prevents phantom images from displaying during the transition
	$('div#rotation div.panel').css({opacity: 0.0});
	$('div#subtextcont span.subtext').css({opacity: 0.0});
	
	//Get the first image and display it (gets set to full opacity)
	$('div#rotation div.panel:first').css({opacity: 1.0});
	$('div#subtextcont span.subtext:first').css({opacity: 1.0});
	$('#thumbs img.thumb').click(function(){switchImage($(this))});
	playRotation()
});

function sizeThumbs() {
	var twidth = $("#thumbs").width();
	var pwidth = $("#thumbcont").width()
	$("#thumbs").css('margin-left',(pwidth - twidth)/2 + 'px')
}

function rotate(count,target) {
	if(target == undefined) {
		var target = 0;
	}
	//Get the first image
	var current = ($('div#rotation div.show')?  $('div#rotation div.show') : $('div#rotation div.panel:first'));
	var currentSubText = ($('div#subtextcont span.show')?  $('div#subtextcont span.show') : $('div#subtextcont span.subtext:first'));
	
	if(target != 0) {
		var next = $("div#panel_" + target);
		var nextSubText = $("span#subtext_" + target);
	} else {
		if(count == 1) {
			//Get next image, when it reaches the end, rotate it back to the first image
			var next = ((current.next().length) ? ((current.next().hasClass('show')) ? $('div#rotation div.panel:first') :current.next()) : $('div#rotation div.panel:first'));
			var nextSubText = ((currentSubText.next().length) ? ((currentSubText.next().hasClass('show')) ? $('div#subtextcont span.subtext:first') :currentSubText.next()) : $('div#subtextcont span.subtext:first'));
		} else {
			var next = ((current.prev().length) ? ((current.prev().hasClass('show')) ? $('div#rotation div.panel:last') :current.prev()) : $('div#rotation div.panel:last'));
			var nextSubText = ((currentSubText.prev().length) ? ((currentSubText.prev().hasClass('show')) ? $('div#subtextcont span.subtext:last') :currentSubText.prev()) : $('div#subtextcont span.subtext:last'));
		}
	}

	
	//Set the fade in effect for the next image, the show class has higher z-index
	nextSubText.css({opacity: 0.0})
	.addClass('show')
	.animate({opacity: 1.0}, 1000);

	next.css({opacity: 0.0})
	.addClass('show')
	.animate({opacity: 1.0}, 1000);


	//Hide the current image
	current.animate({opacity: 0.0}, 1000)
	.removeClass('show');
	currentSubText.animate({opacity: 0.0}, 1000)
	.removeClass('show');

	pauseRotation()
	playRotation()
};

function pauseRotation()
{
	if($('div#rotation div.panel').length > 1) clearInterval(rotateInterval)
}

function playRotation()
{
	if($('div#rotation div.panel').length > 1) rotateInterval = setInterval('rotate(1)',rotationDuration);
}

function switchImage(obj) {
	var id = $(obj).attr('id').split('_')[1]
	rotate(1,id)
}