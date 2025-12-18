# Frontend Structure

This directory contains the frontend files for the GoDir application, organized into separate folders for better maintainability:

## Directory Structure

- `*.html` - Main HTML files
- `css/` - CSS stylesheets
- `js/` - JavaScript files

## Files

### HTML Files
- `index.html` - Main page/home page
- `login.html` - User login/registration page
- `material.html` - Material/file management page

### CSS Files
- `css/index.css` - Styles for the main page
- `css/login.css` - Styles for the login page
- `css/material.css` - Styles for the material page

### JavaScript Files
- `js/index.js` - JavaScript for the main page
- `js/login.js` - JavaScript for the login page
- `js/material.js` - JavaScript for the material page

## Changes Made

Previously, all HTML files contained inline CSS and JavaScript code. These have been separated into external files for better organization and maintainability:

1. Extracted CSS styles from each HTML file into corresponding CSS files in the `css/` directory
2. Extracted JavaScript code from each HTML file into corresponding JS files in the `js/` directory
3. Updated HTML files to reference external CSS and JS files using `<link>` and `<script>` tags

This separation improves:
- Code maintainability
- Reusability
- Readability
- Performance (CSS and JS files can be cached separately)