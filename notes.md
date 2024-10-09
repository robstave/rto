https://jsfiddle.net/hrxcqfyL/7/

// Set up dimensions
const margin = { top: 20, right: 30, bottom: 30, left: 40 }
const width = 800 - margin.left - margin.right
const height = 400 - margin.top - margin.bottom

// Generate date range
const startDate = new Date(2024, 9, 1) // October 1, 2024
const endDate = new Date(2024, 11, 31) // December 31, 2024
const dateRange = d3.timeDay.range(startDate, d3.timeDay.offset(endDate, 1))

// Calculate the number of days
const days = dateRange.length

// Get today's date
const today = new Date()

const svg = d3
  .select("#chart")
  .append("svg")
  .attr("width", width + margin.left + margin.right)
  .attr("height", height + margin.top + margin.bottom)
  .append("g")
  .attr("transform", `translate(${margin.left},${margin.top})`)

// Generate sample data
let data = []
let total = 0
dateRange.forEach((date, index) => {
  const dayOfWeek = date.getDay()
  const isWeekday = dayOfWeek >= 1 && dayOfWeek <= 5
  const comesIn = isWeekday && Math.random() < 0.6 // 60% chance of coming in on weekdays
  total += comesIn ? 7 / days : 0
  data.push({
    date: date,
    comesIn: comesIn,
    total: total,
  })
})

// Set up scales
const xScale = d3.scaleTime().domain([startDate, endDate]).range([0, width])
const yScale = d3
  .scaleLinear()
  .domain([0, d3.max(data, (d) => 4)])
  .range([height, 0])

// Create axes
const xAxis = d3.axisBottom(xScale)
const yAxis = d3.axisLeft(yScale)
svg.append("g").attr("transform", `translate(0,${height})`).call(xAxis)
svg.append("g").call(yAxis)

// Create color zones
const colorScale = d3
  .scaleThreshold()
  .domain([2, 3, 4])
  .range(["red", "yellow", "green"])
const zoneData = [
  { start: 0, end: 2, color: "red", msg:"Oh Dear" },
  { start: 2, end: 3, color: "yellow", msg:"Safe" },
  { start: 3, end: 4, color: "green", msg:"Stellar" },
]
svg
  .selectAll(".zone")
  .data(zoneData)
  .enter()
  .append("rect")
  .attr("class", "zone")
  .attr("x", 0)
  .attr("y", (d) => yScale(d.end))
  .attr("width", width)
  .attr("height", (d) => yScale(d.start) - yScale(d.end))
  .attr("fill", (d) => d.color)
  .attr("opacity", 0.2)

// Ensure the zones are behind the line
svg.selectAll(".zone").lower()

// Add weekend bars
svg.selectAll(".weekend")
  .data(dateRange)
  .enter()
  .filter(d => d.getDay() === 0 || d.getDay() === 6) // 0 is Sunday, 6 is Saturday
  .append("rect")
  .attr("class", "weekend")
  .attr("x", d => xScale(d))
  .attr("y", 0)
  .attr("width", width / days)
  .attr("height", height)
  .attr("fill", "lightgrey")
  .attr("opacity", 0.3)

// Ensure the weekend bars are behind the line
svg.selectAll(".weekend").lower()

// Create line
const line = d3
  .line()
  .x((d) => xScale(d.date))
  .y((d) => yScale(d.total))
svg
  .append("path")
  .datum(data)
  .attr("fill", "none")
  .attr("stroke", "steelblue")
  .attr("stroke-width", 2)
  .attr("d", line)
  
  //---------
  
 

// Add hover line and tooltip
const hoverLine = svg.append("line")
  .attr("class", "hover-line")
  .attr("y1", 0)
  .attr("y2", height)
  .style("stroke", "black")
  .style("stroke-width", "1px")
  .style("stroke-dasharray", "5,5")
  .style("opacity", 0);

const hoverTooltip = d3.select("body").append("div")
  .attr("class", "hover-tooltip")
  .style("opacity", 0)
  .style("position", "absolute")
  .style("background-color", "white")
  .style("border", "solid")
  .style("border-width", "1px")
  .style("border-radius", "5px")
  .style("padding", "10px");

const bisectDate = d3.bisector(d => d.date).left;

svg.append("rect")
  .attr("class", "overlay")
  .attr("width", width)
  .attr("height", height)
  .style("fill", "none")
  .style("pointer-events", "all")
  .on("mouseover", () => {
    hoverLine.style("opacity", 1);
    hoverTooltip.style("opacity", 1);
  })
  .on("mouseout", () => {
    hoverLine.style("opacity", 0);
    hoverTooltip.style("opacity", 0);
  })
  .on("mousemove", mousemove);

function mousemove(event) {
  const [mouseX] = d3.pointer(event);
  const x0 = xScale.invert(mouseX);
  const i = bisectDate(data, x0, 1);
  const d0 = data[i - 1];
  const d1 = data[i];
  const d = x0 - d0.date > d1.date - x0 ? d1 : d0;

  hoverLine
    .attr("x1", xScale(d.date))
    .attr("x2", xScale(d.date));

  hoverTooltip
    .html(`Date: ${d.date.toLocaleDateString()}<br/>Total: ${d.total.toFixed(2)}`)
    .style("left", (event.pageX + 10) + "px")
    .style("top", (event.pageY - 28) + "px");
}
  
  //--

// Add nodes for days you come in
const tooltip = d3.select("body").append("div")
  .attr("class", "tooltip")
  .style("opacity", 0)
  .style("position", "absolute")
  .style("background-color", "white")
  .style("border", "solid")
  .style("border-width", "1px")
  .style("border-radius", "5px")
  .style("padding", "10px")

svg.selectAll(".node")
  .data(data.filter(d => d.comesIn))
  .enter()
  .append("circle")
  .attr("class", "node")
  .attr("cx", d => xScale(d.date))
  .attr("cy", d => yScale(d.total))
  .attr("r", 4)
  .attr("fill", d => d.date > today ? "steelblue" : "darkblue")
  .on("mouseover", (event, d) => {
    tooltip.transition()
      .duration(200)
      .style("opacity", .9);
    tooltip.html(`Date: ${d.date.toLocaleDateString()}<br/>Total: ${d.total.toFixed(2)}`)
      .style("left", (event.pageX + 10) + "px")
      .style("top", (event.pageY - 28) + "px");
  })
  .on("mouseout", (d) => {
    tooltip.transition()
      .duration(500)
      .style("opacity", 0);
  });

// Add title
svg
  .append("text")
  .attr("x", width / 2)
  .attr("y", -margin.top / 2)
  .attr("text-anchor", "middle")
  .style("font-size", "16px")
  .text("Office Days Burnup Chart")
  
  
  
// Add legend
const legend = svg.append("g")
  .attr("class", "legend")
  .attr("transform", `translate(${width-70}, ${height-120})`);

// Color zones legend
zoneData.forEach((zone, i) => {
  legend.append("rect")
    .attr("x", 0)
    .attr("y", i * 20)
    .attr("width", 15)
    .attr("height", 15)
    .attr("fill", zone.color)
    .attr("opacity", 0.2);
  
  legend.append("text")
    .attr("x", 20)
    .attr("y", i * 20 + 12)
    .text(`${zone.msg}`)
    .style("font-size", "12px");
});

// Node color legend
const nodeColors = [
  { color: "darkblue", label: "Past" },
  { color: "steelblue", label: "Future" }
];

nodeColors.forEach((node, i) => {
  legend.append("circle")
    .attr("cx", 7)
    .attr("cy", (i + zoneData.length) * 20 + 7)
    .attr("r", 4)
    .attr("fill", node.color);
  
  legend.append("text")
    .attr("x", 20)
    .attr("y", (i + zoneData.length) * 20 + 12)
    .text(node.label)
    .style("font-size", "12px");
});


// Add vertical line for today
const todayLine = svg.append("line")
  .attr("class", "today-line")
  .attr("x1", xScale(today))
  .attr("y1", 0)
  .attr("x2", xScale(today))
  .attr("y2", height)
  .attr("stroke", "red")
  .attr("stroke-width", 2)
  .attr("stroke-dasharray", "5,5");

// Add label for today's line
svg.append("text")
  .attr("class", "today-label")
  .attr("x", xScale(today))
  .attr("y", 0)
  .attr("dy", "-0.5em")
  .attr("text-anchor", "middle")
  .attr("fill", "red")
  .text("Today");

// Ensure today's line is above the zones but below the data line and points
todayLine.raise();
path.raise();
svg.selectAll(".node").raise();

 