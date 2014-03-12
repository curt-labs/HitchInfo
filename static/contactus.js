function findInvalidFields(id, invalidFields) {
    invalidFields = invalidFields.toLowerCase().split(',');
     
    var formElements = getAllFormElements(document.getElementById(id));
      
    for (var i = 0; i < formElements.length; i++) {       
      if (invalidFields.indexOf(formElements[i].name.toLowerCase()) >= 0 && formElements[i].name != '') {
        formElements[i].className = "invalid";
      }
    }
  }
  
  function findRequiredFields(id, requiredFields) {
    requiredFields = requiredFields.toLowerCase().split(',');
     
    var formElements = getAllFormElements(document.getElementById(id));
      
    for (var i = 0; i < formElements.length; i++) {     
      if (requiredFields.indexOf(formElements[i].name.toLowerCase()) >= 0 && formElements[i].name != '') {
        formElements[i].parentNode.parentNode.parentNode.className = "required";
      }
    }
  }
  
  function getLabelByFor(id) {
    var labels = document.getElementsByTagName('label')
    for (var i = 0; i < labels.length; i++) {
      if (labels[i].htmlFor.toLowerCase() == id.toLowerCase())
        return labels[i];
    }
    return;
  }
  
  function findRequiredFieldLabels(id, requiredFields) {
    requiredFields = requiredFields.toLowerCase().split(',');
    
    var formElements = getAllFormElements(document.getElementById(id));

    for (var i = 0; i < formElements.length; i++) {
      if (requiredFields.indexOf(formElements[i].name.toLowerCase()) >= 0 && formElements[i].name != '') {
        var label =  getLabelByFor(formElements[i].id);
        if (label) {
          //label.innerHTML = label.innerHTML + ' *';
          label.innerHTML = label.innerHTML + ' <img alt="Required Symbol" src="images/reusable/req.gif" title="Required Symbol" />';
        }
      }
    }
  }
    
  // Array.indexOf( value, begin, strict ) - Return index of the first element that matches value
  Array.prototype.indexOf = function( v, b, s ) {
    for( var i = +b || 0, l = this.length; i < l; i++ ) {
      if( this[i]===v || s && this[i]==v ) { return i; }
    }
    return -1;
  };
    
  function getAllFormElements( parent_node ) {
    if( parent_node == undefined ) {
      parent_node = document;
    }
      
    var out = new Array();
    
    formInputs = parent_node.getElementsByTagName("input");      
    for (var i = 0; i < formInputs.length; i++)
    out.push( formInputs.item(i) );
     
    formInputs = parent_node.getElementsByTagName("textarea");
    for (var i = 0; i < formInputs.length; i++)
    out.push( formInputs.item(i) );
      
    formInputs = parent_node.getElementsByTagName("select");
    for (var i = 0; i < formInputs.length; i++)
    out.push( formInputs.item(i) );      
    
    formInputs = parent_node.getElementsByTagName("button");
    for (var i = 0; i < formInputs.length; i++)
    out.push( formInputs.item(i) );    
      
    return out;
  }
  
  function setFormValues(jsonFormStruct) {
     
    //alert(JSON.stringify(jsonFormStruct));
      
    for (var i = 0; i < jsonFormStruct.length; i++) {
      var fieldName = jsonFormStruct[i].FIELDNAME.toLowerCase(); /* struct var names are case sensitive in js */
      var fieldValue = jsonFormStruct[i].FIELDVALUE;
       
      var fieldObject = document.getElementById(fieldName);
       
      if (fieldObject && fieldObject.id == fieldName) {
        if (fieldObject.tagName.toLowerCase() == 'input') {
          if (fieldObject.type == 'text' || fieldObject.type == 'hidden') fieldObject.value = fieldValue;
        }
        else if (fieldObject.tagName.toLowerCase() == 'textarea')
          fieldObject.innerHTML = fieldValue;
      }
      else {
        var fieldObjects = document.getElementsByName(fieldName);
        
        for (var j = 0; j < fieldObjects.length; j++) {
          if (fieldObjects[j].tagName.toLowerCase() == 'input') {
            if (fieldObjects[j].type == 'checkbox' || fieldObjects[j].type == 'radio')
              if (fieldValue.indexOf(fieldObjects[j].value) >= 0) fieldObjects[j].checked = "true";
          }
        }
      }
    }
  }