# http://www.logarithmic.net/pfh/blog/0118662041


class Searcher:
    def __init__(self, filename):
        self.f = open(filename, 'rb')
        self.f.seek(0,2)
        self.length = self.f.tell()

    def find(self, string):
        low = 0
        high = self.length
        while low < high:
            mid = (low+high)//2
            p = mid
            while p >= 0:
                self.f.seek(p)
                if self.f.read(1) == '\n': break
                p -= 1
            if p < 0: self.f.seek(0)
            line = self.f.readline()
            print '--', mid, line
            if line < string:
                low = mid+1
            else:
                high = mid

        p = low
        while p >= 0:
            self.f.seek(p)
            if self.f.read(1) == '\n': break
            p -= 1
        if p < 0: self.f.seek(0)

        result = [ ]
        while True:
            line = self.f.readline()
            if not line or not line.startswith(string): break
            if line[-1:] == '\n': line = line[:-1]
            result.append(line[len(string):])
        return result


bs = Searcher("sample.txt")

print bs.find("00000006Hello@World")