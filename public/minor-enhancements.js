(function() {
  function createElement(type, parent) {
    const elem = document.createElement(type)
    parent.appendChild(elem)
    return elem
  }

  const elements = document.getElementsByClassName("field-row")
  for (var i = 0; i < elements.length; i++) {
    const element = elements[i]
    const elementTitle = element.firstElementChild.firstElementChild.innerHTML.split(":")[0]
    const removeButtonTd = createElement("td", element)
    const removeButton = createElement("button", removeButtonTd)
    removeButton.onclick = function() {
      if (confirm("Are you sure you want to delete the field '" + elementTitle + "'?")) {
        element.remove()
      }
    }
    removeButton.innerHTML = "Remove field"
  }
})()
