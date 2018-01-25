#!/bin/sh
#
awk -F'|' 'BEGIN {
        printf("digraph G {\n")
        printf("\tfontname = \"Bitstream Vera Sans\"\n")
        printf("\tfontsize = 8\n")
        printf("\n")
        printf("\tgraph [\n")
        #printf("\t\tnodesep = 6\n")
        printf("\t\tranksep = 1\n")
        printf("\t]\n")
        printf("\n")
        printf("\tnode [\n")
        printf("\t\tfontname = \"Bitstream Vera Sans\"\n")
        printf("\t\tfontsize = 8\n")
        printf("\t\tshape = record\n")
        printf("\t\tstyle = filled\n")
        printf("\t]\n")
        printf("\n")
        printf("\tedge [\n")
        printf("\t\tfontname = \"Bitstream Vera Sans\"\n")
        printf("\t\tfontsize = 8\n")
        printf("\t]\n")

        group = ""
        class = ""
        edges = ""
}
{
        if ((start > 0) && (NF > 1)) {
                a = $1; gsub(/[ ]*/, "", a)
                b = $2; gsub(/[ ]*/, "", b)
                c = $3; gsub(/[ ]*/, "", c)
                label = sprintf("%s%s : %s %s\\l", label, a, b, c)
                if (b == "") {
                        edges = sprintf("%s\t%s -> %s\n", edges, tag, a)
                }
        }
}
/^---/ {
        start = 1
}
/<div/ {
        start = 0
}
/^$/ {
        start = 0
}
/^# / {
        if (class ne "") {
                printf("\t\t\tlabel = \"{%s||%s}\"\n", class, label)
                printf("\t\t]\n")
        }
        if (group ne "") {
                printf("\t}\n")
        }
        group = $0; gsub(/#[ ]*/, "", group)
        tag = group; gsub(/[ ]*/, "", tag)
        class = ""

        printf("\tsubgraph cluster%s {\n", tag)
        printf("\t\tlabel = \"%s\"\n", group)
}
/^## / {
        if (class ne "") {
                printf("\t\t\tlabel = \"{%s||%s}\"\n", class, label)
                printf("\t\t]\n")
        }
        class = $0; gsub(/##[ ]*/, "", class)
        tag = class; gsub(/[ ]*/, "", tag)
        printf("\t\t%s [\n", tag)
        label = ""
}
END {
        if (class ne "") {
                printf("\t\t\tlabel = \"{%s||%s}\"\n", class, label)
                printf("\t\t]\n")
        }
        if (group ne "") {
                printf("\t}\n")
        }


        printf("\tedge [\n")
        printf("\t\tarrowhead = \"empty\"\n")
        printf("\t]\n")
        printf("\n")
        printf("%s\n", edges)
        printf("\n")
        printf("}\n")
}' README.md
