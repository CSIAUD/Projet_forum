// https://observablehq.com/@d3/radial-stacked-bar-chart-ii@171
export default function define(runtime, observer) {
  const main = runtime.module();
  const fileAttachments = new Map([["data-2.csv",new URL("./files/33954e9770a615014df5bcb8be3e257700ec99c4eb9a3206db60577f0e63bf1f217621ab5d7eefa3b187ecb5418db4609bef452ae939370a9d89a709f803b31f",import.meta.url)]]);
  main.builtin("FileAttachment", runtime.fileAttachments(name => fileAttachments.get(name)));
  main.variable(observer()).define(["md"], function(md){return(
md`# Radial Stacked Bar Chart II

Radial layouts are pretty, but may impede comparison, so consider them primarily to emphasize cyclical patterns. (This example’s data is not cyclical, so the radial layout here is gratuitous!)`
)});
  main.variable(observer("chart")).define("chart", ["d3","DOM","width","height","data","z","arc","xAxis","yAxis","legend"], function(d3,DOM,width,height,data,z,arc,xAxis,yAxis,legend)
{
  const svg = d3.select(DOM.svg(width, height))
      .attr("viewBox", `${-width*1.5 /2} ${-height*1.5/2} ${width*1.5 } ${height*1.5}`)
      .style("width", "100%")
      .style("height", "auto")
      .style("font", "15px sans-serif");

  svg.append("g")
    .selectAll("g")
    .data(d3.stack().keys(data.columns.slice(1))(data))
    .join("g")
      .attr("fill", d => z(d.key))
    .selectAll("path")
    .data(d => d)
    .join("path")
      .attr("d", arc);

  svg.append("g")
      .call(xAxis);

  svg.append("g")
      .call(yAxis);

  svg.append("g")
      .call(legend);

  return svg.node();
}
);
  main.variable(observer("data")).define("data", ["d3","FileAttachment"], async function(d3,FileAttachment){return(
d3.csvParse(await FileAttachment("data-2.csv").text(), (d, _, columns) => {
  let total = 0;
  for (let i = 1; i < columns.length; ++i) total += d[columns[i]] = +d[columns[i]];
  d.total = total;
  return d;
}).sort((a, b) => b.total - a.total)
)});
  main.variable(observer("arc")).define("arc", ["d3","y","x","innerRadius"], function(d3,y,x,innerRadius){return(
d3.arc()
    .innerRadius(d => y(d[0]))
    .outerRadius(d => y(d[1]))
    .startAngle(d => x(d.data.Category))
    .endAngle(d => x(d.data.Category) + x.bandwidth())
    .padAngle(0.01)
    .padRadius(innerRadius)
)});
  main.variable(observer("x")).define("x", ["d3","data"], function(d3,data){return(
d3.scaleBand()
    .domain(data.map(d => d.Category))
    .range([0, 2 * Math.PI])
    .align(0)
)});
  main.variable(observer("y")).define("y", ["d3","data","innerRadius","outerRadius"], function(d3,data,innerRadius,outerRadius){return(
d3.scaleRadial()
    .domain([0, d3.max(data, d => d.total)])
    .range([innerRadius, outerRadius])
)});
  main.variable(observer("z")).define("z", ["d3","data"], function(d3,data){return(
d3.scaleOrdinal()
    .domain(data.columns.slice(1))
    .range(["#98abc5", "#8a89a6", "#7b6888", "#6b486b", "#a05d56", "#d0743c", "#ff8c00"])
)});
  main.variable(observer("xAxis")).define("xAxis", ["data","x","innerRadius"], function(data,x,innerRadius){return(
g => g
    .attr("text-anchor", "middle")
    .call(g => g.selectAll("g")
      .data(data)
      .enter().append("g")
        .attr("transform", d => `
          rotate(${((x(d.Category) + x.bandwidth() / 2) * 180 / Math.PI - 90)})
          translate(${innerRadius},0)
        `)
        .call(g => g.append("line")
            .attr("x2", -5)
            .attr("stroke", "#000"))
        .call(g => g.append("text")
            .attr("transform", d => (x(d.Category) + x.bandwidth() / 2 + Math.PI / 2) % (2 * Math.PI) < Math.PI
                ? "rotate(90) translate(0,16)"
                : "rotate(-90) translate(0,-9)")
            .text(d => d.Category)))
)});
  main.variable(observer("yAxis")).define("yAxis", ["y"], function(y){return(
g => g
    .attr("text-anchor", "end")
    .call(g => g.append("text")
        .attr("x", -6)
        .attr("y", d => -y(y.ticks(10).pop()))
        .attr("dy", "-1em")
        .text("Population"))
    .call(g => g.selectAll("g")
      .data(y.ticks(10).slice(1))
      .join("g")
        .attr("fill", "none")
        .call(g => g.append("circle")
            .attr("stroke", "#000")
            .attr("stroke-opacity", 0.5)
            .attr("r", y))
        .call(g => g.append("text")
            .attr("x", -6)
            .attr("y", d => -y(d))
            .attr("dy", "0.35em")
            .attr("stroke", "#fff")
            .attr("stroke-width", 5)
            .text(y.tickFormat(10, "s"))
         .clone(true)
            .attr("fill", "#000")
            .attr("stroke", "none")))
)});
  main.variable(observer("legend")).define("legend", ["data","z"], function(data,z){return(
g => g.append("g")
  .selectAll("g")
  .data(data.columns.slice(1).reverse())
  .join("g")
    .attr("transform", (d, i) => `translate(-40,${(i - (data.columns.length - 1) / 2) * 20})`)
    .call(g => g.append("rect")
        .attr("width", 18)
        .attr("height", 18)
        .attr("fill", z))
    .call(g => g.append("text")
        .attr("x", 24)
        .attr("y", 9)
        .attr("dy", "0.35em")
        .text(d => d))
)});
  main.variable(observer("width")).define("width", function(){return(
1000
)});
  main.variable(observer("height")).define("height", ["width"], function(width){return(
width
)});
  main.variable(observer("innerRadius")).define("innerRadius", function(){return(
180
)});
  main.variable(observer("outerRadius")).define("outerRadius", ["width","height"], function(width,height){return(
Math.min(width, height) * 0.67
)});
  main.variable(observer("d3")).define("d3", ["require"], function(require){return(
require("d3@6")
)});
  return main;
}
