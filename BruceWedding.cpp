#include <iostream>
#include <fstream>
#include <string>
#include <algorithm>
 
// This should be the maximum distance expected between
// data points. For example, if you're largest value is 10000 (5 chars)
// And your file looks like 9998,9999,10000, then your number should be 6
// accounting for the comma
const int MaxDataLen = 10;
 
std::string& ltrim(std::string& str, const std::string& chars = "\t\n\v\f\r ")
{
	str.erase(0, str.find_first_of(chars)+1);
	return str;
}
 
int GetValueAt(std::ifstream &myfile, unsigned long midpoint)
{
	char buffer[MaxDataLen * 3] = { 0 };
	// We seek to the point and subtract the max data len to
	// ensure we get all the data at that point. Note, it is 
	// possible to backup to the next data point. We can check
	// for that by counting the commas in our sample
	myfile.seekg(midpoint - MaxDataLen);
	myfile.read(buffer, MaxDataLen * 2 - 2);
	std::string myString = buffer;
	// Trim off the first comma found
	// because that's where our value begins in most cases
	myString = ltrim(myString, ",");
 
	// If there are more than 1 commas left, trim it until there aren't
	while (std::count(myString.begin(), myString.end(), ',') > 1)
	{
		myString = ltrim(myString, ",");
	}
	// Note: The string will still have another comma in it, but stoi will take care of discarding it
	return std::stoi(myString);
}
 
int binarySearch(std::ifstream &myfile, int key, int& location)
{
	std::streampos end;
	myfile.seekg(0, std::ios::end);
	end = myfile.tellg();
 
	int fSize = end;
	int low = 0, high = fSize - 1, midpoint = 0, curValue;
	while (low <= high)
	{
		midpoint = low + (high - low) / 2;
 
		curValue = GetValueAt(myfile, midpoint);
		if (key == curValue)
		{
			location = midpoint;
			return location;
		}
		else 
			if (key < curValue)
				high = midpoint - 1;
		else
			low = midpoint + 1; 
	}
	return location;
}
 
int main() 
{
	std::ifstream myfile("c:/temp/example.txt");
	int location = -1;  // not found	
	location = binarySearch(myfile, 5000, location);
	std::cout << "5000 was found at location: " << location << std::endl;
 
	location = binarySearch(myfile, 50, location);
	std::cout << "50 was found at location: " << location << std::endl;
 
	location = binarySearch(myfile, 500000, location);
	std::cout << "500,000 was found at location: " << location << std::endl;
 
	location = binarySearch(myfile, 314159, location);
	std::cout << "314159 was found at location: " << location << std::endl;
 
	myfile.close();
	std::cout << "Done!";
	return 0;
}
