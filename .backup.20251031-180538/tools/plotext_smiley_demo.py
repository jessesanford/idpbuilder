#!/usr/bin/env python3
import sys

try:
    from textual.app import App, ComposeResult
    from textual.widget import Widget
    from textual.widgets import Static
    from textual.containers import Vertical
    from textual_plotext import PlotextPlot  # type: ignore
except Exception as e:
    print("This demo requires textual and textual-plotext. Install with: pip install textual textual-plotext", file=sys.stderr)
    sys.exit(1)


class Smiley(Widget):
    def compose(self) -> ComposeResult:
        yield PlotextPlot(id="plot")

    def on_mount(self) -> None:  # type: ignore[override]
        plot = self.query_one(PlotextPlot)
        plt = plot.plt
        try:
            plt.clear_figure()
        except Exception:
            pass
        # Eyes
        plt.scatter([-1, 1], [1, 1])
        # Mouth (arc)
        mx = [-1.2, -0.6, 0.0, 0.6, 1.2]
        my = [-0.2, -0.5, -0.8, -0.5, -0.2]
        plt.plot(mx, my)
        # Ensure there is at least one text-like content by plotting a tiny point at center
        plt.scatter([0.0], [0.0])
        try:
            plt.xlim(-2, 2)
            plt.ylim(-2, 2)
            plt.canvas_color("black")
        except Exception:
            pass
        # Avoid plt.text() on versions where it triggers a signals/Text mismatch
        try:
            plot.refresh()
        except Exception:
            pass


class PlotextSmileyApp(App):
    CSS = """
    Screen { layout: vertical; }
    #title { height: 1; }
    #body { height: 1fr; }
    """

    def compose(self) -> ComposeResult:
        yield Static("Plotext Smiley Demo (press Ctrl+C to exit)", id="title")
        yield Smiley(id="body")


if __name__ == "__main__":
    app = PlotextSmileyApp()
    app.run()
