#!/usr/bin/env python3
import matplotlib.pyplot as plt
import numpy as np
import re
import seaborn as sns

# Set the style using seaborn
sns.set_theme(style="whitegrid")

# Read benchmark data
with open('benchmark.txt', 'r') as f:
    data = f.read()

# Parse benchmark results
def parse_benchmarks(data):
    results = {}
    for line in data.split('\n'):
        if line.startswith('Benchmark'):
            parts = line.split('\t')
            if len(parts) >= 3:
                # Split the benchmark name and handle the -8 suffix
                name_parts = parts[0].split('/')
                if len(name_parts) >= 2:
                    test_type = name_parts[0].replace('Benchmark', '')
                    text_type = name_parts[1].split('-')[0]  # Remove the -8 suffix
                    time = float(re.search(r'(\d+) ns/op', parts[2]).group(1))
                    
                    if test_type not in results:
                        results[test_type] = {'time': {}}
                    
                    if text_type not in results[test_type]['time']:
                        results[test_type]['time'][text_type] = []
                    
                    results[test_type]['time'][text_type].append(time)
    
    return results

results = parse_benchmarks(data)

# Calculate averages
def get_averages(results):
    averages = {}
    for test_type in results:
        averages[test_type] = {'time': {}}
        for text_type in results[test_type]['time']:
            values = results[test_type]['time'][text_type]
            averages[test_type]['time'][text_type] = np.mean(values)
    return averages

averages = get_averages(results)

# Prepare data for plotting
text_types = list(next(iter(averages.values()))['time'].keys())
num_tests = len(text_types)

# Create figure with extra space at top for title
fig = plt.figure(figsize=(10, 3.5*num_tests + 1))

# Add main title at top left with margin
title_ax = fig.add_axes([0, 0.96, 1, 0.04])
title_ax.axis('off')
title_ax.text(0.02, 0.5, 'Tokenization Performance Comparison', 
             fontsize=16, fontweight='bold', ha='left', va='center')
title_ax.text(0.02, 0, '(Lower is Better)', 
             fontsize=12, ha='left', va='top', style='italic')

# Create subplot grid for the charts with spacing
gs = fig.add_gridspec(num_tests, 1, top=0.92, hspace=0.4)
axes = [fig.add_subplot(gs[i]) for i in range(num_tests)]

# Colors for different implementations
colors = {'CL100kTokenizer': '#2ecc71', 'TiktokenCL100k': '#e74c3c',
         'O200kTokenizer': '#3498db', 'TiktokenO200k': '#e67e22'}

for idx, text_type in enumerate(text_types):
    ax = axes[idx]
    
    # Get data for this text type
    data = []
    labels = []
    colors_list = []
    
    # Group implementations by type
    impl_pairs = [
        ('CL100k', ['CL100kTokenizer', 'TiktokenCL100k']),
        ('O200k', ['O200kTokenizer', 'TiktokenO200k'])
    ]
    
    for impl_type, impls in impl_pairs:
        for impl in impls:
            if impl in averages:
                data.append(averages[impl]['time'][text_type])
                labels.append(f'{impl_type}: {impl}')
                colors_list.append(colors[impl])
    
    # Create horizontal bars
    y_pos = np.arange(len(data))
    bars = ax.barh(y_pos, data)
    
    # Color the bars
    for bar, color in zip(bars, colors_list):
        bar.set_color(color)
    
    # Add percentage differences
    for i in range(0, len(data), 2):
        if i+1 < len(data):
            base = data[i]
            comp = data[i+1]
            diff_percent = ((comp - base) / base) * 100
            ax.text(comp + (comp * 0.02), i+1, f'+{diff_percent:.0f}%', va='center')
    
    # Customize the plot
    ax.set_yticks(y_pos)
    ax.set_yticklabels(labels)
    ax.set_title(f'{text_type} Text', pad=10, fontsize=12)
    ax.set_xlabel('Time (ns/op)')
    
    # Add value labels on the bars
    for i, v in enumerate(data):
        ax.text(v/2, i, f'{int(v):,}', 
                ha='center', va='center', color='white', fontweight='bold')

plt.tight_layout()
plt.savefig('results.png', dpi=300, bbox_inches='tight')
