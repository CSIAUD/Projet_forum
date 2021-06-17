// https://observablehq.com/@d3/pannable-chart@265
export default function define(runtime, observer) {
  const content = (document.getElementById("allData").innerText).replaceAll("|", "\n")
  const main = runtime.module();
  const fileAttachments = new Map([["aapl.csv",new URL("./files/de259092d525c13bd10926eaf7add45b15f2771a8b39bc541a5bba1e0206add4880eb1d876be8df469328a85243b7d813a91feb8cc4966de582dc02e5f8609b7",import.meta.url)]]);
  main.builtin("FileAttachment", runtime.fileAttachments(date => fileAttachments.get(date)));
  main.variable(observer()).define(["md"], function(md){return(
md`# Pannable Chart

This [area chart](/@d3/area-chart) supports horizontal panning. Try scrolling left or right below.`
)});
  main.variable(observer("chart")).define("chart", ["x","data","margin","d3","width","height","yAxis","xAxis","area"], function*(x,data,margin,d3,width,height,yAxis,xAxis,area)
{
  const minX = x(data[0].date);
  const maxX = x(data[data.length - 1].date);
  const overwidth = maxX - minX + margin.left + margin.right;

  const parent = d3.create("div");

  parent.append("svg")
      .attr("width", width)
      .attr("height", height)
      .style("position", "absolute")
      .style("pointer-events", "none")
      .style("z-index", 1)
      .call(svg => svg.append("g").call(yAxis));

  const body = parent.append("div")
      .style("overflow-x", "scroll")
      .style("-webkit-overflow-scrolling", "touch");

  body.append("svg")
      .attr("width", overwidth)
      .attr("height", height)
      .style("display", "block")
      .call(svg => svg.append("g").call(xAxis))
    .append("path")
      .datum(data)
      .attr("fill", "rgb(252,141,89)")
      .attr("d", area);

  yield parent.node();

  // Initialize the scroll offset after yielding the chart to the DOM.
  body.node().scrollBy(overwidth, 0);
}
);
  main.variable(observer("height")).define("height", function(){return(
420
)});
  main.variable(observer("margin")).define("margin", function(){return(
{top: 20, right: 20, bottom: 30, left: 40}
)});
  main.variable(observer("x")).define("x", ["d3","data","margin","width"], function(d3,data,margin,width){return(
d3.scaleUtc()
    .domain(d3.extent(data, d => d.date))
    .range([margin.left, width * 6 - margin.right])
)});
  main.variable(observer("y")).define("y", ["d3","data","height","margin"], function(d3,data,height,margin){return(
d3.scaleLinear()
    .domain([0, d3.max(data, d => d.value)]).nice(6)
    .range([height - margin.bottom, margin.top])
)});
  main.variable(observer("xAxis")).define("xAxis", ["height","margin","d3","x","width"], function(height,margin,d3,x,width){return(
g => g
    .attr("transform", `translate(10,${height - margin.bottom})`)
    .call(d3.axisBottom(x).ticks(d3.utcMonth.every(1200 / width)).tickSizeOuter(0))
)});
  main.variable(observer("yAxis")).define("yAxis", ["margin","d3","y","data"], function(margin,d3,y,data){return(
g => g
    .attr("transform", `translate(${margin.left},0)`)
    .call(d3.axisLeft(y).ticks(6))
    .call(g => g.select(".domain").remove())
    .call(g => g.select(".tick:last-of-type text").clone()
        .attr("x", 3)
        .attr("text-anchor", "start")
        .attr("font-weight", "bold")
        .text(data.y))
)});
  main.variable(observer("area")).define("area", ["d3","x","y"], function(d3,x,y){return(
d3.area()
    .curve(d3.curveStep)
    .x(d => x(d.date))
    .y0(y(0))
    .y1(d => y(d.value))
)});
  main.variable(observer("data")).define("data", ["d3","FileAttachment"], async function(d3,FileAttachment){return(
Object.assign(d3.csvParse(content, d3.autoType).map(({date, nb}) => ({date, value: nb})), {y: "Nb Post/Jour"})
)});
  main.variable(observer("d3")).define("d3", ["require"], function(require){return(
require("d3@6")
)});
  return main;
}
