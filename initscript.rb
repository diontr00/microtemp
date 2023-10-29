#!/usr/bin/env ruby

require 'fileutils'
# Check if a template name argument is provided
if ARGV.length != 1
  puts "Usage: #{$PROGRAM_NAME} <template_name>"
  exit 1
end

# Capture the template name from the command-line argument
template_name = ARGV[0]

# Function to replace template in a file
def replace_template(file, template_name)
  # Read the file content
  content = File.read(file)

  # Replace all occurrences of {{{mytemplate}}} with the provided template name
  updated_content = content.gsub('{{{mytemplate}}}', template_name)

  # Write the updated content back to the file
  File.open(file, 'w') { |f| f.write(updated_content) }
end

# Recursively find all files in the current directory and its subdirectories, excluding .git
def find_and_replace_files(directory, template_name)
  Dir.foreach(directory) do |item|
    next if item == '.' || item == '..' || item == '.git'

    path = File.join(directory, item)

    if File.file?(path)
      replace_template(path, template_name)
    elsif File.directory?(path)
      find_and_replace_files(path, template_name)
    end
  end
end

# Start the search and replacement from the current directory
find_and_replace_files('.', template_name)

 "Template replacement complete."

system("go mod init #{template_name}")
system("go mod tidy")



FileUtils.rm_rf('.git')

system('git init -b main')

FileUtils.rm(__FILE__)

