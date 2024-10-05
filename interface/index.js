import * as d3 from "https://cdn.jsdelivr.net/npm/d3@7/+esm";

const maxZoom = 8,
  minZoom = 0.5;
let dragging = false;

const clearStates = () => {
  d3.selectAll(".node,.edge").classed(
    "clicked ignored mentioned parent child incoming outgoing animated",
    false
  );
};

d3.json("reachability.json").then((reachability) => {
  d3.xml("dot.svg").then((data) => {
    // append svg data
    const svgContainer = d3.select("div#svgContainer");
    svgContainer.node().append(data.documentElement);
    svgContainer.on("dblclick", () => clearStates());

    const appBoundingClientRect = d3
      .select("div#app")
      .node()
      .getBoundingClientRect();

    const getViewBox = (svg) =>
      svg
        .attr("viewBox")
        .split(" ")
        .map((n) => parseFloat(n));

    const originalViewBox = getViewBox(d3.select("svg"));

    const nodeEventToScreenCoord = (event) => {
      const svg = d3.select("svg").node();
      // get the svg coordinate
      const [svgX, svgY] = d3.pointer(event, svg);
      let svgPoint = new DOMPoint(svgX, svgY);
      // transform svg coordinate to screen coordinate
      return svgPoint.matrixTransform(svg.getScreenCTM());
    };

    const getReachableNodes = (id) =>
      Object.entries(reachability[id].nodeReachability)
        .filter(([key, value]) => value > 0)
        .map(([key, value]) => key);

    function handleSvgMousedown(event) {
      event.stopPropagation();
      event.preventDefault();
      dragging = true;
      d3.select(this).style("cursor", "grab");
    }

    function handleSvgMouseup(event) {
      event.stopPropagation();
      event.preventDefault();
      dragging = false;
      d3.select(this).style("cursor", "default");
    }
    // set dragging on svg
    function handleSvgDrag(event) {
      if (dragging) {
        event.stopPropagation();
        event.preventDefault();
        d3.select("#resetBtn").style("opacity", 1);
        const svg = d3.select(this);

        const [startX, startY, startWidth, startHeight] = getViewBox(svg);

        const startClient = {
          x: event.clientX,
          y: event.clientY,
        };

        let newSVGPoint = new DOMPoint(startClient.x, startClient.y);
        let CTM = svg.node().getScreenCTM();
        const startSVGPoint = newSVGPoint.matrixTransform(CTM.inverse());

        let moveToClient = {
          x: event.clientX + event.movementX,
          y: event.clientY + event.movementY,
        };

        newSVGPoint = new DOMPoint(moveToClient.x, moveToClient.y);
        CTM = svg.node().getScreenCTM();
        const moveToSVGPoint = newSVGPoint.matrixTransform(CTM.inverse());

        const delta = {
          x: startSVGPoint.x - moveToSVGPoint.x,
          y: startSVGPoint.y - moveToSVGPoint.y,
        };

        svg.attr(
          "viewBox",
          `${startX + delta.x} ${startY + delta.y} ${startWidth} ${startHeight}`
        );
      }
    }
    // set zoom on svg
    function handleSvgZoom(event) {
      event.stopPropagation();
      event.preventDefault();
      d3.select("#resetBtn").style("opacity", 1);
      const svg = d3.select(this);

      const [startX, startY, startWidth, startHeight] = getViewBox(svg);

      const startClient = {
        x: event.clientX,
        y: event.clientY,
      };

      const newSVGPoint = new DOMPoint(startClient.x, startClient.y);
      let CTM = svg.node().getScreenCTM();
      const startSVGPoint = newSVGPoint.matrixTransform(CTM.inverse());

      const r = event.deltaY > 0 ? 1.1 : event.deltaY < 0 ? 0.9 : 1;

      const zoomedWidth = startWidth * r,
        zoomedHeight = startHeight * r;

      const zoom = originalViewBox[2] / zoomedWidth;
      if (zoom >= maxZoom || zoom <= minZoom) return;

      svg.attr("viewBox", `${startX} ${startY} ${zoomedWidth} ${zoomedHeight}`);

      CTM = svg.node().getScreenCTM();
      const moveToSVGPoint = newSVGPoint.matrixTransform(CTM.inverse());

      const delta = {
        x: startSVGPoint.x - moveToSVGPoint.x,
        y: startSVGPoint.y - moveToSVGPoint.y,
      };

      const [
        intermediateX,
        intermediateY,
        intermediateWidth,
        intermediateHeight,
      ] = getViewBox(svg);
      svg.attr(
        "viewBox",
        `${intermediateX + delta.x} ${
          intermediateY + delta.y
        } ${intermediateWidth} ${intermediateHeight}`
      );
    }

    // create a tooltip
    const tooltip = d3
      .select("div#svgContainer")
      .append("div")
      .style("display", "none")
      .style("opacity", 0)
      .classed("tooltip", true);

    // Three function that change the tooltip when user hover / move / leave a cell
    function handleNodeMouseover(event) {
      event.preventDefault();
      tooltip
        .style("opacity", d3.select(this).classed("ignored") ? 0.5 : 1) // ignored node's tooltip should be more transparent
        .style("display", "flex");
    }
    function handleNodeMousemove(event) {
      event.preventDefault();
      const nodeId = d3.select(this).node().id;
      const nodeProps = reachability[nodeId];
      const screenCoord = nodeEventToScreenCoord(event);

      // in case tooltip renders out of the viewport
      const xAlign =
        screenCoord.x < appBoundingClientRect.right - 150 ? "left" : "right";
      const yAlign =
        screenCoord.y < appBoundingClientRect.bottom - 150 ? "top" : "bottom";
      const x =
        xAlign === "right"
          ? appBoundingClientRect.right - screenCoord.x
          : screenCoord.x;
      const y =
        yAlign === "bottom"
          ? appBoundingClientRect.bottom - screenCoord.y
          : screenCoord.y;

      tooltip
        .html(
          `
        <h4>${nodeId}</h4>
            <span class="tooltip-prop">
                <span class="prop-key">
                    Parents
                </span>
                <span class="prop-val">
                    ${nodeProps.parents ? nodeProps.parents.length : 0}
                </span>
            </span>
            <span class="tooltip-prop">
                <span class="prop-key">
                    Children
                </span>
                <span class="prop-val">
                    ${nodeProps.children ? nodeProps.children.length : 0}
                </span>
            </span>
            <span class="tooltip-prop">
                <span class="prop-key">
                    Offsprings
                </span>
                <span class="prop-val">
                    ${getReachableNodes(nodeId).length}
                </span>
            </span>
        </div>
        `
        )
        .style("top", null)
        .style("bottom", null)
        .style("left", null)
        .style("right", null)
        .style(xAlign, x + 10 + "px")
        .style(yAlign, y + 10 + "px");
    }
    function handleNodeMouseleave(event) {
      event.preventDefault();
      tooltip.style("opacity", 0).style("display", "none");
    }
    function handleNodeClick(event) {
      event.stopPropagation();
      event.preventDefault();
      const clickedNode = this;
      const clicked = d3.select(clickedNode).classed("clicked");

      clearStates();

      if (!clicked) {
        const shouldAnimate = d3.select("#animationToggle").node().checked;
        d3.select(clickedNode)
          .classed("clicked", true)
          .classed("animated", shouldAnimate);

        const reachableNodes = getReachableNodes(clickedNode.id);

        // classify all reachable nodes
        d3.selectAll(".node").each(function (p, j) {
          d3.select(this)
            .classed(
              "mentioned",
              reachableNodes && reachableNodes.includes(this.id)
            )
            .classed(
              "parent",
              reachability[clickedNode.id].parents &&
                reachability[clickedNode.id].parents.includes(this.id)
            )
            .classed(
              "child",
              reachability[clickedNode.id].children &&
                reachability[clickedNode.id].children.includes(this.id)
            )
            .classed("animated", shouldAnimate);
        });

        // classify all reachable edges
        d3.selectAll(".edge").each(function (p, j) {
          d3.select(this)
            .classed(
              "mentioned",
              reachability[clickedNode.id].reachableEdges &&
                reachability[clickedNode.id].reachableEdges.includes(this.id)
            )
            .classed(
              "incoming",
              reachability[clickedNode.id].incomingEdges &&
                reachability[clickedNode.id].incomingEdges.includes(this.id)
            )
            .classed(
              "outgoing",
              reachability[clickedNode.id].outgoingEdges &&
                reachability[clickedNode.id].outgoingEdges.includes(this.id)
            )
            .classed("animated", shouldAnimate);
        });

        d3.selectAll(
          ".node:not(#" +
            clickedNode.id +
            "):not(.mentioned):not(.parent):not(.child)"
        ).classed("ignored", true);
        d3.selectAll(
          ".edge:not(.mentioned):not(.incoming):not(.outgoing)"
        ).classed("ignored", true);
      }
    }

    // bundle svg events
    d3.select("svg")
      .on("wheel", handleSvgZoom)
      .on("mousedown", handleSvgMousedown)
      .on("mouseup", handleSvgMouseup)
      .on("mousemove", handleSvgDrag);

    // bundle node events
    d3.selectAll(".node")
      .on("mouseover", handleNodeMouseover)
      .on("mousemove", handleNodeMousemove)
      .on("mouseleave", handleNodeMouseleave)
      .on("dblclick", handleNodeClick);

    d3.select("#resetBtn").on("click", function(event) {
      d3.select("svg").attr("viewBox", originalViewBox.join(" "));
      d3.select(this).style("opacity", 0);
    });
  });
});
