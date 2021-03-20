滑动窗口算法(Sliding Window Algorithm)是常见的一种算法：它思想简洁且功能强大，可以用来解决一些查找满足一定条件的连续区间的性质/长度的问题。由于区间连续，因此当区间发生变化时，可以通过旧有的计算结果对搜索空间进行剪枝，这样便减少了重复计算，降低了时间复杂度，它还可以将嵌套的循环问题，转换为单循环问题，同样也是降低时间复杂度。