import matplotlib.pyplot as plt
import json
import numpy as np
from matplotlib.patches import Wedge
import matplotlib.patches as mpatches

def read_data():
    with open("maps.json") as file:
        return json.load(file)
    return None

def main():
    maps = read_data()
    for data in maps["data"]:
            # Create figure and axis
        fig, ax = plt.subplots(figsize=(10, 12))

# FOV parameters
        fov_radius = 2.0  # Longer radius as requested
        fov_angle_deg = 5  # +/- 5 degrees
        fov_angle_rad = np.deg2rad(fov_angle_deg)

# Set axis limits first
        ax.set_xlim(-0.2, 1.2)
        ax.set_ylim(-1.3, 1.3)

# Plot FOV for player 1 pieces (drawn first so they're behind)
        for piece in data['Player1Pieces']:
            angle_start = np.rad2deg(piece['Angle'] - fov_angle_rad)
            angle_end = np.rad2deg(piece['Angle'] + fov_angle_rad)
            wedge = Wedge((piece['Position']['x'], piece['Position']['y']), 
                        fov_radius, angle_start, angle_end, 
                        facecolor='gray', alpha=0.2, edgecolor='none')
            ax.add_patch(wedge)
            wedge.set_clip_box(ax.bbox)

# Plot FOV for player 2 pieces
        for piece in data['Player2Pieces']:
            angle_start = np.rad2deg(piece['Angle'] - fov_angle_rad)
            angle_end = np.rad2deg(piece['Angle'] + fov_angle_rad)
            wedge = Wedge((piece['Position']['x'], piece['Position']['y']), 
                        fov_radius, angle_start, angle_end, 
                        facecolor='gray', alpha=0.2, edgecolor='none')
            ax.add_patch(wedge)
            wedge.set_clip_box(ax.bbox)

# Plot player 1 segments (dark lines)
        for seg in data['player1Segments']:
            ax.plot([seg['p1']['x'], seg['p2']['x']], 
                    [seg['p1']['y'], seg['p2']['y']], 
                    'k-', linewidth=2)

# Plot player 2 segments (dark lines)
        for seg in data['player2Segments']:
            ax.plot([seg['p1']['x'], seg['p2']['x']], 
                    [seg['p1']['y'], seg['p2']['y']], 
                    'k-', linewidth=2)

# Plot player 1 pieces (blue circles)
        for piece in data['Player1Pieces']:
            ax.plot(piece['Position']['x'], piece['Position']['y'], 
                    'bo', markersize=8)

# Plot player 2 pieces (red circles)
        for piece in data['Player2Pieces']:
            ax.plot(piece['Position']['x'], piece['Position']['y'], 
                    'ro', markersize=8)

# Plot player 1 king (larger blue circle)
        ax.plot(data['player1King']['x'], data['player1King']['y'], 
                'bo', markersize=12)

# Plot player 2 king (larger red circle)
        ax.plot(data['player2King']['x'], data['player2King']['y'], 
                'ro', markersize=12)

# Draw the two boundary squares with gray thin lines
# Square 1: min x = 0, max x = 1, min y = 0.1, max y = 1.1
        square1 = [[0, 0.1], [1, 0.1], [1, 1.1], [0, 1.1], [0, 0.1]]
        ax.plot([p[0] for p in square1], [p[1] for p in square1], 
                'gray', linewidth=1, linestyle='-')

# Square 2: min x = 0, max x = 1, min y = -0.1, max y = -1.1
        square2 = [[0, -0.1], [1, -0.1], [1, -1.1], [0, -1.1], [0, -0.1]]
        ax.plot([p[0] for p in square2], [p[1] for p in square2], 
                'gray', linewidth=1, linestyle='-')

# Set equal aspect ratio and labels
        ax.set_aspect('equal')
        ax.set_xlabel('X')
        ax.set_ylabel('Y')
        ax.set_title('Game State Visualization')
        ax.grid(True, alpha=0.3)

# Apply the clipping to x: [0, 1] and y: [-1.1, 1.1]
        ax.set_xlim(0, 1)
        ax.set_ylim(-1.1, 1.1)

        plt.tight_layout()
        plt.show()

main()
