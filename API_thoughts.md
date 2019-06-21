# Thoughts on basic API

Event-driven?
Poll-driven?
Async? Sync?

I think the following backends would be nice:
  * tcell (non-graphics)
  * ebiten (graphical potential)

As such, we would need Interface(s) for:
  "Window" (holding window)
    "Buffer" (target for render)
      "Cell"

Probably each "Buffer" should be able to contain other buffers? Would this be the basis for multiple types of UI Elements, such as DropdownList, etc.?

I think it'd be nice to use go's channels for much of the functionality.

Perhaps:

```
gr.UseBackend(gr.EBITEN)
gr.UseBackend(gr.TCELL)

for {
  <-gr.hasEvent

  select {
    case window := <-gr.windowChannel:
      handleWindow(window)
    case key := <-gr.keyChannel:
      handleKey(key)
    case mouse := <-gr.mouseChannel:
      handleMouse(mouse)
  }
}

or

window, err := gr.NewWindow()

window.Init() // This triggers ebiten.Run(gr.update...)
window.EnableMouse()
window.Clear()
window.SetColumns(80)
window.SetRows(24)
// window.Size() will be defaults.

for event := gr.WaitEvent() {
  switch event := ev.(type) {
    case *gr.EventWindow:
    case *gr.EventKey:
    case *gr.EventMouse:
    case *gr.EventDraw:
  }
}

//

for event := <-gr.EventChan {
  switch event := event.(type) {
    case *gr.EventWindow:
      // event.window.Clear()...
    case *gr.EventKey:
    case *gr.EventMouse:
    case *gr.EventDraw:

  }
}

```
