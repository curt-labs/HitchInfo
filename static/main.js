function addEvent(obj, evType, fn){
  if (obj.addEventListener){
    obj.addEventListener(evType, fn, true);
    return true;
  } else if (obj.attachEvent){
  var r = obj.attachEvent("on"+evType, fn);
  return r;
  } else {
    return false;
  }
}

function updateModifyType(){

	var makeField = "selectVehicleMake2";
	var yearField = "selectVehicleYear2";
	var modelField = "selectVehicleModel2";
	var mountField = "selectVehicleMount2";

	var modifyType = document.getElementById("modifyType");

	if(modifyType.checked){
		document.getElementById(mountField).disabled = false;
		document.getElementById(makeField).disabled = false;
		document.getElementById(yearField).disabled = false;
		document.getElementById(modelField).disabled = false;
	}else{
		document.getElementById(mountField).disabled = true;
		document.getElementById(makeField).disabled = true;
		document.getElementById(yearField).disabled = true;
		document.getElementById(modelField).disabled = true;
	}

}

function getYears2(){

	var makeField = "selectVehicleMake2";
	var yearField = "selectVehicleYear2";
	var modelField = "selectVehicleModel2";
	var mountField = "selectVehicleMount2";

	// create logging function
	var myOnLog = function(msg){
	}

	// create completion function
	var myOnComplete = function(responseText, responseXML){

		//document.getElementById('outputText').value = responseText;
		var modelNodes = responseXML.getElementsByTagName('vehicle-year');
		var modelNode = '';
		var formObject = document.getElementById(yearField);

		var i = 0;
		var formOption= '';
		var yearName = '';
		var yearValue = '';

		//remove the make
		document.getElementById(makeField).options.length=0;

		//go through and add all the current makes
		formObject.options.length=0;
		formObject.options.length=modelNodes.length;
		for(i=0;i<modelNodes.length;i++){
			modelNode = modelNodes[i];
			yearName = modelNode.getElementsByTagName('year')[0].firstChild.nodeValue;
			yearValue = modelNode.getElementsByTagName('id')[0].firstChild.nodeValue;
			formOption = new Option(yearName,yearValue);
			formObject.options[i] = formOption;
		}

		//refresh the models
		getMakes2();


	}

	// create provider instance; wire events
	var provider = new oyXMLRPCProvider();
	provider.onComplete = myOnComplete;
	provider.onLog = myOnLog;
	provider.onError = myOnLog;

	var mountid = document.getElementById(mountField).options[document.getElementById(mountField).selectedIndex].value;

	provider.submit("index.cfm?event=yearxml&mount=" + mountid);
}

function getMakes2(){

	var makeField = "selectVehicleMake2";
	var yearField = "selectVehicleYear2";
	var modelField = "selectVehicleModel2";
	var mountField = "selectVehicleMount2";

	// create logging function
	var myOnLog = function(msg){
	}

	// create completion function
	var myOnComplete = function(responseText, responseXML){

		//document.getElementById('outputText').value = responseText;
		var modelNodes = responseXML.getElementsByTagName('vehicle-make');
		var modelNode = '';
		var formObject = document.getElementById(makeField);

		var i = 0;
		var formOption= '';
		var makeName = '';
		var makeValue = '';

		//remove the model
		document.getElementById(modelField).options.length=0;

		//go through and add all the current makes
		formObject.options.length=0;
		formObject.options.length=modelNodes.length;
		for(i=0;i<modelNodes.length;i++){
			modelNode = modelNodes[i];
			makeName = modelNode.getElementsByTagName('make')[0].firstChild.nodeValue;
			makeValue = modelNode.getElementsByTagName('id')[0].firstChild.nodeValue;
			formOption = new Option(makeName,makeValue);
			formObject.options[i] = formOption;
		}

		//refresh the models
		getModels2();


	}

	// create provider instance; wire events
	var provider = new oyXMLRPCProvider();
	provider.onComplete = myOnComplete;
	provider.onLog = myOnLog;
	provider.onError = myOnLog;

	var yearId = document.getElementById(yearField).options[document.getElementById(yearField).selectedIndex].value;
	var mountid = document.getElementById(mountField).options[document.getElementById(mountField).selectedIndex].value;

	provider.submit("index.cfm?event=makexml&year=" + yearId + "&mount=" + mountid);
}

function getModels2(){

	var makeField = "selectVehicleMake2";
	var yearField = "selectVehicleYear2";
	var modelField = "selectVehicleModel2";
	var styleField = "selectVehicleStyle2";
	var mountField = "selectVehicleMount2";

	// create logging function
	var myOnLog = function(msg){
	}

	// create completion function
	var myOnComplete = function(responseText, responseXML){

		//document.getElementById('outputText').value = responseText;
		var modelNodes = responseXML.getElementsByTagName('model');
		var modelNode = '';
		var formObject = document.getElementById(modelField);

		var i = 0;
		var formOption= '';
		var modelName = '';
		var modelValue = '';

		//remove the model
		document.getElementById(styleField).options.length=0;

		//go through and add all the current models
		formObject.options.length=0;
		formObject.options.length=modelNodes.length;
		for(i=0;i<modelNodes.length;i++){
			modelNode = modelNodes[i];
			modelName = modelNode.getElementsByTagName('name')[0].firstChild.nodeValue;
			modelValue = modelNode.getElementsByTagName('id')[0].firstChild.nodeValue;
			formOption = new Option(modelName,modelValue);
			formObject.options[i] = formOption;
		}

		//refresh the models
		getStyles2();


	}

	// create provider instance; wire events
	var provider = new oyXMLRPCProvider();
	provider.onComplete = myOnComplete;
	provider.onLog = myOnLog;
	provider.onError = myOnLog;

	var yearId = document.getElementById(yearField).options[document.getElementById(yearField).selectedIndex].value;
	var makeId = document.getElementById(makeField).options[document.getElementById(makeField).selectedIndex].value;
	var mountid = document.getElementById(mountField).options[document.getElementById(mountField).selectedIndex].value;

	provider.submit("index.cfm?event=modelxml&year=" + yearId + "&make=" + makeId + "&mount=" + mountid);
}

function getStyles2(){

	var makeField = "selectVehicleMake2";
	var yearField = "selectVehicleYear2";
	var modelField = "selectVehicleModel2";
	var styleField = "selectVehicleStyle2";
	var mountField = "selectVehicleMount2";

	// create logging function
	var myOnLog = function(msg){
	}

	// create completion function
	var myOnComplete = function(responseText, responseXML){
		var vehicleTypeNodes = responseXML.getElementsByTagName('vehicle-style');
		var vehicleTypeNode = '';
		var formObject = document.getElementById(styleField);

		var i = 0;
		var formOption= '';
		var vehicleTypeModel = '';
		var vehicleTypeValue = '';


		//go through and add all the yearsresponseXML
		formObject.options.length=0;
		formObject.options.length=vehicleTypeNodes.length;
		for(i=0;i<vehicleTypeNodes.length;i++){
			vehicleTypeNode = vehicleTypeNodes[i];
			vehicleTypeModel = vehicleTypeNode.getElementsByTagName('vstyle')[0].firstChild.nodeValue;
			vehicleTypeValue = vehicleTypeNode.getElementsByTagName('id')[0].firstChild.nodeValue;
			formOption = new Option(vehicleTypeModel,vehicleTypeValue);
			formObject.options[i] = formOption;
		}



	}


	// create provider instance; wire events
	var provider = new oyXMLRPCProvider();
	provider.onComplete = myOnComplete;
	provider.onLog = myOnLog;
	provider.onError = myOnLog;

	var yearId = document.getElementById(yearField).options[document.getElementById(yearField).selectedIndex].value;
	var makeId = document.getElementById(makeField).options[document.getElementById(makeField).selectedIndex].value;
	var modelId = document.getElementById(modelField).options[document.getElementById(modelField).selectedIndex].value;
	var mountid = document.getElementById(mountField).options[document.getElementById(mountField).selectedIndex].value;
	provider.submit("index.cfm?event=stylexml&year=" + yearId + "&make=" + makeId + "&model=" + modelId + "&mount=" + mountid);


}





function isValid(type, str) {
  if (type.toLowerCase() === "email") {
    var reEmail = new RegExp(/^(("[\w-\s]+")|([\w-]+(?:\.[\w-]+)*)|("[\w-\s]+")([\w-]+(?:\.[\w-]+)*))(@((?:[\w-]+\.)*\w[\w-]{0,66})\.([a-z]{2,6}(?:\.[a-z]{2})?)$)|(@\[?((25[0-5]\.|2[0-4][0-9]\.|1[0-9]{2}\.|[0-9]{1,2}\.))((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[0-9]{1,2})\.){2}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[0-9]{1,2})\]?$)/i);
    if (!reEmail.test(str)) return false;
    return true;
  } else if (type.toLowerCase() === "telephone") {
    var rePhoneNumber = new RegExp(/^(\+\d)*\s*(\(?\d{3}\)?\s*)*-?\d{3}(-{0,1}|\s{0,1})\d{2}(-{0,1}|\s{0,1})\d{2}$/);
    if (!rePhoneNumber.test(str)) return false;
    return true;
  }
};

jQuery.fn.simpleAccordion = function() {
  return this.each(function() {
    var height = 0;
    $(this).find("ul").each(function() {
      if($(this).height() > height) height = $(this).height();
    });
    $(this).find("ul").height(height).hide();
    if($(this).find("li.active").size() > 0) $(this).find("li.active ul").show()
    else $(this).find("ul:first").show();
    var $$ = $(this);
    $(this).find("li a").click(
      function() {
        var checkElement = $(this).next();
        if((checkElement.is("ul")) && (checkElement.is(":visible"))) {
          return false;
        }
        if((checkElement.is("ul")) && (!checkElement.is(":visible"))) {
          $$.find("li.active").removeClass("active");
          $$.find("ul:visible").slideUp("normal");
          checkElement.parent("li").addClass("active");
          checkElement.slideDown("normal");
          return false;
        }
      }
    );
  });
};

(function($) {
  $(function() {
    $("input[type=text][title]").each(function() { $(this).val($(this).attr("title")); if($.trim($(this).val()) == "") $(this).val($(this).attr("title")); $(this).focus(function() { if($(this).val() == $(this).attr("title")) $(this).val(""); }).blur(function() { if($.trim($(this).val()) == "") $(this).val($(this).attr("title")); }); });
    $("a[href][rel*=external]").attr("target", "_blank");
    $(".fade-hover").hoverIntent(function() { $(this).fadeTo("fast", 0.5); }, function() { $(this).fadeTo("fast", 1.0); });
    $(".lo").hoverIntent(function() { $(this).removeClass("lo").addClass("hi"); }, function() { $(this).removeClass("hi").addClass("lo"); });


    $(".glossary").each(function() { $(this).attr("rel", "shadowbox;width=400;height=400;"); $(this).attr("href", "index.cfm?event=glossary.term&term=" + $(this).text()); $(this).attr("title", $(this).text()); });

    $('a[href=http://www.buycurthitches.com/]').closest('tr').remove();

    Shadowbox.init();

    $("#faq-menu").simpleAccordion();
  });
})(jQuery);