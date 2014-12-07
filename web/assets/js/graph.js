var arrMaxLen = 10;

var graph = new joint.dia.Graph;

var paper = new joint.dia.Paper({
    el: $('#paper'),
    width: 800,
    height: 600,
    gridSize: 1,
    model: graph
});

function newPoint(x, y) {
    return {
        x: x,
        y: y
    }
}

function clearGraph() {
    graph.clear()
}

function drawItem(point, item) {
    if (item.type == "undefined") {
        return
    }

    switch (item.type) {
        case "block":
            drawBlock(point, item);
            break;
        case "assignment":
            drawAssignment(point, item);
            break;
    }
}

function drawBlock(point, block) {
    if (block.children == "undefined" || block.children.length == 0) {
        return
    }
    var children = block.children

    for (k in children) {
        if (children[k] != null) {
            drawItem(point, children[k])
        }
    }
}

function drawAssignment(point, stmt) {
    if (stmt.left == "undefined" || stmt.right == "undefined") {
        return
    }

    drawValue(point, stmt.right)
}

function drawValue(point, val) {
    if (val.type == "undefined") {
        return
    }

    switch (val.type) {
        case "array":
            drawArray(point, val.items)
        case "const":
            drawConst(point, val)
    }
}

function drawArray(point, arr) {
    for (i = 0;
        (i < arr.length) && (i < arrMaxLen); i++) {

        drawValue(point, arr[i])
    }
}

function drawConst(point, val) {
    cellH = 20
    cellW = 20
    var cell = new joint.shapes.basic.Rect({
        position: {
            x: point.x,
            y: point.y
        },
        size: {
            width: cellW,
            height: cellH
        },
        attrs: {
            text: {
                text: val.value
            }
        }
    });

    point.x += cellW
    point.y += cellH

    graph.addCell(cell);
}

function state(point, label) {

    var cell = new joint.shapes.fsa.State({
        position: {
            x: x,
            y: y
        },
        size: {
            width: 60,
            height: 60
        },
        attrs: {
            text: {
                text: label
            }
        }
    });
    graph.addCell(cell);
    return cell;
};
