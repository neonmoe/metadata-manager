(function() {
  function createElement(type, parent) {
    const elem = document.createElement(type)
    parent.appendChild(elem)
    return elem
  }

  function createButton(text, className, parent, func) {
    const removeButton = createElement("button", parent)
    removeButton.classList.add(className)
    removeButton.type = "button"
    removeButton.innerHTML = text
    removeButton.onclick = func
  }

  function moveElement(element, indexDiff) {
    const parent = element.parentNode
    const newList = []
    let insertIndex = indexDiff
    let elementIterationIndex = 0
    for (let i = 0; i < parent.childNodes.length; i++) {
      if (parent.childNodes[i].nodeType != Node.ELEMENT_NODE) {
        continue
      }

      if (element == parent.childNodes[i]) {
        insertIndex += newList.length
      } else  {
        newList.push(parent.childNodes[i])
      }
    }

    if (insertIndex < 0) insertIndex += newList.length
    if (insertIndex >= newList.length) insertIndex -= newList.length

    newList.splice(insertIndex, 0, element)
    for (let i = 0; i < newList.length; i++) {
      parent.appendChild(newList[i])
    }
  }

  const elements = document.getElementsByClassName("field-row")
  for (var i = 0; i < elements.length; i++) {
    const element = elements[i]
    const elementTitle = element.firstElementChild.firstElementChild.innerHTML.split(":")[0]
    const td = createElement("td", element)

    createButton("×", "remove-button", td, function() {
      if (confirm("Are you sure you want to delete the field '" + elementTitle + "'?")) {
        element.remove()
      }
    })

    createButton("↑", "move-button", td, function() {
      moveElement(this.parentNode.parentNode, -1)
    })

    createButton("↓", "move-button", td, function() {
      moveElement(this.parentNode.parentNode, 1)
    })
  }
})()
