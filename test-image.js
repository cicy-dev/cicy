const blessed = require('blessed');
const fs = require('fs');
const path = require('path');

// Create screen
const screen = blessed.screen({
    smartCSR: true,
    title: 'image-test',
    fullUnicode: true
});

// Simple test image (red dot)
const testImageBase64 = 'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==';

// Save to file
const tempFile = '/tmp/test_image.png';
fs.writeFileSync(tempFile, Buffer.from(testImageBase64, 'base64'));

// Create image widget
const image = blessed.image({
    parent: screen,
    top: 1,
    left: 1,
    width: 20,
    height: 10,
    file: tempFile,
    type: 'ansi'  // Use ANSI art
});

// Message
const msg = blessed.box({
    parent: screen,
    top: 12,
    left: 1,
    width: 50,
    height: 3,
    content: 'Image test complete!',
    style: {
        fg: 'green'
    }
});

screen.key(['q', 'C-c'], () => {
    process.exit(0);
});

screen.render();
