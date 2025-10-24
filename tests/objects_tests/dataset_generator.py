import os
import json
import random
from types import coroutine
from shapely.geometry import Point, LineString
from shapely.lib import intersects

FILE_NAMES = [
    "point.json",
    "line.json",
    "segment.json"
]
TESTS_PER_METHOD = 100



# ---                     ---
# --- 1. Helper Functions ---
# ---                     ---

# --- Initializers ---
def make_point():
    return Point(random.uniform(-10, 10), random.uniform(-10, 10))

def make_segment():
    return LineString([make_point(), make_point()])

def make_segment_paralel_to(l):
    coords = list(l.coords)
    a = (coords[1][1] - coords[0][1])/(coords[1][0] - coords[0][0])
    b = random.uniform(-10, 10)
    x1 = random.uniform(-10, 10)
    x2 = random.uniform(-10, 10)
    p1 = Point(x1, a * x1 + b)
    p2 = Point( x2, a * x2 + b)
    return LineString([p1, p2])

def make_line():
    a = random.uniform(-100, 100)
    b = random.uniform(-10, 10)
    p1 = Point(-100000, a * -100000 + b)
    p2 = Point( 100000, a * 100000 + b)
    return LineString([p1, p2])

def make_line_parallel_to(l):
    coords = list(l.coords)
    a = (coords[1][1] - coords[0][1])/(coords[1][0] - coords[0][0])
    b = random.uniform(-10, 10)
    p1 = Point(-100000, a * -100000 + b)
    p2 = Point( 100000, a * 100000 + b)
    return LineString([p1, p2])


# --- To json ---
def point_to_json(p):
    return {"x": float_to_json(p.x), "y": float_to_json(p.y)}

def line_to_json(s):
    coords = list(s.coords)
    return {
        "p1": point_to_json(Point(coords[0])),
        "p2": point_to_json(Point(coords[1])),
    }

def segment_to_json(s):
    coords = list(s.coords)
    return {
        "p1": point_to_json(Point(coords[0])),
        "p2": point_to_json(Point(coords[1])),
    }

def float_to_json(f):
    return round(f, 6)


# ---                    ---
# --- 3. Test Generators ---
# ---                    ---
def create_point_test_dataset():
    test_json = {}

    # Test dataset to test distance method
    distance_dataset = []
    for _ in range(TESTS_PER_METHOD):
        p1 = make_point()
        p2 = make_point()
        distance = p1.distance(p2)
        test_case = (point_to_json(p1), point_to_json(p2), float_to_json(distance))
        distance_dataset.append(test_case)
    test_json["distance"] = distance_dataset

    return test_json

def create_line_test_dataset():
    test_json = {}

    # Test dataset to test shortest distance method
    distance_dataset = []
    for _ in range(TESTS_PER_METHOD):
        p = make_point()
        l = make_line()
        distance = l.distance(p)
        test_case = (point_to_json(p), line_to_json(l), float_to_json(distance))
        distance_dataset.append(test_case)
    test_json["distance"] = distance_dataset

    # Test dataset to test intersection with line method
    intersection_dataset = []
    for _ in range(TESTS_PER_METHOD - 10):
        test_case = ()
        l1 = make_line()
        l2 = make_line()
        if l1.intersects(l2):
            point = l1.intersection(l2).centroid
            test_case = (line_to_json(l1), line_to_json(l2), point_to_json(point))
        else:
            test_case = (line_to_json(l1), line_to_json(l2), None)
        intersection_dataset.append(test_case)
    # Special test case: parralel lines
    for _ in range(10):
        l1 = make_line()
        l2 = make_line_parallel_to(l1)
        if l1.intersects(l2):
            raise Exception("Error, parralel lines should never intersect!")
        else:
            test_case = (line_to_json(l1), line_to_json(l2), None)
        intersection_dataset.append(test_case)
    test_json["lineIntersection"] = intersection_dataset

    # Test dataset to test intersection with segment method
    intersection_dataset = []
    for _ in range(TESTS_PER_METHOD - 10):
        test_case = ()
        l = make_line()
        s = make_segment()
        if l.intersects(s):
            point = l.intersection(s).centroid
            test_case = (line_to_json(l), segment_to_json(s), point_to_json(point))
        else:
            test_case = (line_to_json(l), segment_to_json(s), None)
        intersection_dataset.append(test_case)
    # Special test case: paralel Segment 
    for _ in range(10):
        l = make_line()
        s = make_segment_paralel_to(l)
        if l.intersects(s):
            raise Exception("Error, parralel line and segment should never intersect!")
        else:
            test_case = (line_to_json(l), segment_to_json(s), None)
        intersection_dataset.append(test_case)

    test_json["segmentIntersection"] = intersection_dataset

    return test_json

def create_segment_test_dataset():
    test_json = {}

    # Test dataset to test shortest distance method
    distance_dataset = []
    for _ in range(TESTS_PER_METHOD):
        p = make_point()
        s = make_segment()
        distance = s.distance(p)
        test_case = (segment_to_json(s), point_to_json(p), float_to_json(distance))
        distance_dataset.append(test_case)
    test_json["distance"] = distance_dataset

    # Test dataset to test intersection with segment method
    intersection_dataset = []
    for _ in range(TESTS_PER_METHOD - 10):
        s1 = make_segment()
        s2 = make_segment()
        if s1.intersects(s2):
            point = s1.intersection(s2).centroid
            test_case = (segment_to_json(s1), segment_to_json(s2), point_to_json(point))
        else:
            test_case = (segment_to_json(s1), segment_to_json(s2), None)
        intersection_dataset.append(test_case)
    for _ in range(10):
        s1 = make_segment()
        s2 = make_segment_paralel_to(s1)
        if s1.intersects(s2):
            raise Exception("Error, parralel line and segment should never intersect!")
        else:
            test_case = (segment_to_json(s1), segment_to_json(s2), None)
        intersection_dataset.append(test_case)
    test_json["segmentIntersection"] = intersection_dataset

    return test_json

# ---                            ---
# --- 4. Putting it all together ---
# ---                            ---
def generate_test_datasets():
    for file_name in FILE_NAMES:
        if not os.path.isfile(file_name):
            under_test = file_name.split('.')[0]
            if under_test == "point":
                test_json = create_point_test_dataset()
                with open(f"{file_name}", "w") as f:
                    json.dump(test_json, f, indent=4)
            elif under_test == "line":
                test_json = create_line_test_dataset()
                with open(f"{file_name}", "w") as f:
                    json.dump(test_json, f, indent=4)
            elif under_test == "segment":
                test_json = create_segment_test_dataset()
                with open(f"{file_name}", "w") as f:
                    json.dump(test_json, f, indent=4)

if __name__ == "__main__":
    generate_test_datasets()
