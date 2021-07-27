OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;
bit[2]   c;

reset q;

h  q[0];
cx q[0], q[1];

c[0] = measure q[0];
c[1] = measure q[1];
