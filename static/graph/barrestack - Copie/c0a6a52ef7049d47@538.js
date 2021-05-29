// https://observablehq.com/@d3/stacked-bar-chart@538
import define1 from "./a33468b95d0b15b0@703.js";

export default function define(runtime, observer) {
  const main = runtime.module();
  const fileAttachments = new Map([["us-population-state-age.csv",new URL("./files/cacf3b872e296fd3cf25b9b8762dc0c3aa1863857ecba3f23e8da269c584a4cea9db2b5d390b103c7b386586a1104ce33e17eee81b5cc04ee86929f1ee599bfe",import.meta.url)]]);
  main.builtin("FileAttachment", runtime.fileAttachments(Category => fileAttachments.get(Category)));
  main.variable(observer()).define(["md"]);
  main.variable(observer("key")).define("key", ["legend","color"], function(legend,color){return(
legend({color, title: "Role"})
)});
  main.variable(observer("chart")).define("chart", ["d3","width","height","series","color","x","y","formatValue","xAxis","yAxis"], function(d3,width,height,series,color,x,y,formatValue,xAxis,yAxis)
{
  const svg = d3.create("svg")
      .attr("viewBox", [0, 0, width, height]);

  svg.append("g")
    .selectAll("g")
    .data(series)
    .join("g")
      .attr("fill", d => color(d.key))
    .selectAll("rect")
    .data(d => d)
    .join("rect")
      .attr("x", (d, i) => x(d.data.Category))
      .attr("y", d => y(d[1]))
      .attr("height", d => y(d[0]) - y(d[1]))
      .attr("width", x.bandwidth())
    .append("title")
      .text(d => `${d.data.Category} ${d.key}
${formatValue(d.data[d.key])}`);

  svg.append("g")
      .call(xAxis);

  svg.append("g")
      .call(yAxis);

  return svg.node();
}
);
  main.variable(observer("data")).define("data", ["d3","FileAttachment"], async function(d3,FileAttachment){return(
d3.csvParse(await FileAttachment("us-population-state-age.csv").text(), (d, i, columns) => (d3.autoType(d), d.total = d3.sum(columns, c => d[c]), d)).sort((a, b) => b.total - a.total)
)});
  main.variable(observer("series")).define("series", ["d3","data"], function(d3,data){return(
d3.stack()
    .keys(data.columns.slice(1))
  (data)
    .map(d => (d.forEach(v => v.key = d.key), d))
)});
  main.variable(observer("x")).define("x", ["d3","data","margin","width"], function(d3,data,margin,width){return(
d3.scaleBand()
    .domain(data.map(d => d.Category))
    .range([margin.left, width - margin.right])
    .padding(0.1)
)});
  main.variable(observer("y")).define("y", ["d3","series","height","margin"], function(d3,series,height,margin){return(
d3.scaleLinear()
    .domain([0, d3.max(series, d => d3.max(d, d => d[1]))])
    .rangeRound([height - margin.bottom, margin.top])
)});
  main.variable(observer("color")).define("color", ["d3","series"], function(d3,series){return(
d3.scaleOrdinal()
    .domain(series.map(d => d.key))
    .range(d3.schemeSpectral[series.length])
    .unknown("#ccc")
)});
  main.variable(observer("xAxis")).define("xAxis", ["height","margin","d3","x"], function(height,margin,d3,x){return(
g => g
    .attr("transform", `translate(0,${height - margin.bottom})`)
    .call(d3.axisBottom(x).tickSizeOuter(0))
    .call(g => g.selectAll(".domain").remove())
)});
  main.variable(observer("yAxis")).define("yAxis", ["margin","d3","y"], function(margin,d3,y){return(
g => g
    .attr("transform", `translate(${margin.left},0)`)
    .call(d3.axisLeft(y).ticks(null, "s"))
    .call(g => g.selectAll(".domain").remove())
)});
  main.variable(observer("formatValue")).define("formatValue", function(){return(
x => isNaN(x) ? "N/A" : x.toLocaleString("en")
)});
  main.variable(observer("height")).define("height", function(){return(
600
)});
  main.variable(observer("margin")).define("margin", function(){return(
{top: 10, right: 10, bottom: 20, left: 40}
)});
  main.variable(observer("d3")).define("d3", ["require"], function(require){return(
require("d3@6")
)});
  const child1 = runtime.module(define1);
  main.import("legend", child1);
  return main;
}
