import json
import matplotlib.pyplot as plt

def read_data():
    with open("maps.json") as file:
        return json.load(file)
    return None

def plot_map(mapData):
    # Player 1 map:
    plt.figure()

    for item in mapData["player1Segments"]:
        plt.plot([item["p1"]["x"], item["p2"]["x"]], [item["p1"]["y"], item["p2"]["y"]], 'b-', linewidth=2)

    plt.plot(mapData["player1King"]["x"], mapData["player1King"]["y"], 'ro', label='King')

    plt.xlim(0, 1)
    plt.ylim(0, 1)

    plt.grid(True)
    plt.waitforbuttonpress()
    plt.close()

    # Player 2 map:
    plt.figure()

    for item in mapData["player2Segments"]:
        plt.plot([item["p1"]["x"], item["p2"]["x"]], [item["p1"]["y"], item["p2"]["y"]], 'b-', linewidth=2)

    plt.plot(mapData["player2King"]["x"], mapData["player2King"]["y"], 'ro', label='King')

    plt.grid(True)
    plt.waitforbuttonpress()
    plt.close()

def plot_maps(data):
    for item in data["data"]:
        plot_map(item)

def main():
    data = read_data()
    plot_maps(data)

main()
