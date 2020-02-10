/**
 * Find the position of the start of the first line in the file that is
 * greater than or equal to the target line, using a binary search.
 * 
 * @param file
 *            the file to search.
 * @param target
 *            the target to find.
 * @return the position of the first line that is greater than or equal to
 *         the target line.
 * @throws IOException
 */
public static long search(RandomAccessFile file, String target)
        throws IOException {
    /*
     * because we read the second line after each seek there is no way the
     * binary search will find the first line, so check it first.
     */
    file.seek(0);
    String line = file.readLine();
    if (line == null || line.compareTo(target) >= 0) {
        /*
         * the start is greater than or equal to the target, so it is what
         * we are looking for.
         */
        return 0;
    }

    /*
     * set up the binary search.
     */
    long beg = 0;
    long end = file.length();
    while (beg <= end) {
        /*
         * find the mid point.
         */
        long mid = beg + (end - beg) / 2;
        file.seek(mid);
        file.readLine();
        line = file.readLine();

        if (line == null || line.compareTo(target) >= 0) {
            /*
             * what we found is greater than or equal to the target, so look
             * before it.
             */
            end = mid - 1;
        } else {
            /*
             * otherwise, look after it.
             */
            beg = mid + 1;
        }
    }

    /*
     * The search falls through when the range is narrowed to nothing.
     */
    file.seek(beg);
    file.readLine();
    return file.getFilePointer();
