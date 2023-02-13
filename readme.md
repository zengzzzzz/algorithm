1. Skip List / Quick List
A skip list is a data structure that is similar to an ordered linked list but with some additional features. It uses a probabilistic approach to achieve an average O(log n) time complexity for search, insert, and delete operations, where n is the number of elements in the list.

A quick list is a specific implementation of a skip list where each node in the list is a linked list of nodes with the same key value. This allows for efficient handling of duplicates and improves the performance of certain operations.

Skip lists were first described by William Pugh in his 1990 paper "Skip Lists: A Probabilistic Alternative to Balanced Trees".

Pros
Simple and easy to understand.
Fast average time complexity for search, insert, and delete operations.
Can handle duplicates efficiently.
Cons
Not as efficient as balanced trees in the worst case.
Space overhead due to the additional pointers used in the skip list structure.

2. B-Tree
A B-Tree is a type of balanced tree data structure that is commonly used in databases, file systems, and other applications where fast access to large amounts of ordered data is required.

B-Trees are optimized for operations on blocks of data stored on disk, such as disk drives or flash drives. They work by keeping the keys and data in a node closely packed and maintaining a large number of keys in each node, reducing the number of disk accesses required to find a key.

B-Trees are also self-balancing, meaning that they automatically adjust the structure of the tree to maintain balance and ensure that the height of the tree remains logar to the number of keys stored in the tree. This ensures that search, insert, and delete operations have a logarithmic time complexity, even in the worst case.

Pros
Self-balancing, ensuring logarithmic time complexity for search, insert, and delete operations in the worst case.
Optimized for operations on blocks of data stored on disk.
Efficient storage of large amounts of ordered data.
Cons
More complex than simple linked lists or skip lists.
Can have a larger space overhead due to the additional information stored in each node.

3. B+ Tree
A B+ Tree is a variant of the B-Tree data structure that is commonly used in databases, file systems, and other applications where fast access to large amounts of ordered data is required.

Like B-Trees, B+ Trees are self-balancing and optimized for operations on blocks of data stored on disk. However, B+ Trees differ from B-Trees in that all data is stored in the leaves of the tree, with the internal nodes only containing keys. This can lead to more efficient storage and retrieval of data in certain situations.

B+ Trees also support range queries more efficiently than B-Trees, as all data is stored in the leaves and can be easily retrieved by following a path from the root to the appropriate leaf node.

Pros
Self-balancing, ensuring logarithmic time complexity for search, insert, and delete operations in the worst case.
Optimized for operations on blocks of data stored on disk.
Supports range queries more efficiently than B-Trees.
Efficient storage and retrieval of data in certain situations.
Cons
More complex than simple linked lists or skip lists.
Can have a larger space overhead due to the additional information stored in each node.
Note: It's important to choose the appropriate data structure for your specific use case based on the size and type of data you will be working with, as well as the operations you will be performing on that data.

3.1 limit count

4 reverse binary iteration （redis scan）

5 raft

6 radix tree

7 quorum

8 radix tree

9 zab / vsr / paxos

10 perim / kruskal

11 npc

12 lru

13 random sampling without replacement
