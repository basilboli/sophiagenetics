import React  from 'react';
import { scaleLinear } from 'd3-scale';
import { Axis, axisPropsFromTickScale, LEFT} from 'react-d3-axis';
import './App.css';

const HistogramChart = (props) => {

    var width = 250;
    var height = 100;
    
    // let's use max value to properly scale y axis
    let yDomain = [0, props.max];

    let yScale = scaleLinear()
        .range([height, 0])
        .domain(yDomain);

    let setBarHeight = d => (height - yScale(d))
    let setRectY = d => (yScale(d))

    var transformBar = (d, i) => ("translate(" + i * width / props.data.length + ",0)");
    return (
        <svg>
            <g transform="translate(40,20)">
                <text x="20" y="-10" className="App-chart-header">{props.title}</text>
                {props.data.map((d, i) => {
                    return <g key={i} className="bar" transform={transformBar(d, i)} >
                        <rect width={width / props.data.length} height={setBarHeight(d)} y={setRectY(d)} fill="#bae8e8" />
                        <text x="0" y={height+20} class="App-axis">{i+1}</text>
                    </g>
                })}
                
            </g>
            <g transform="translate(40,20)">
                <Axis {...axisPropsFromTickScale(yScale, 10)} style={{ orient: LEFT }} />
            </g>
        </svg>
    )
}

export default HistogramChart;