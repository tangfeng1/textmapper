/**
 * Copyright 2002-2013 Evgeny Gryaznov
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.textmapper.templates.types.ast;

import org.textmapper.templates.types.TypesTree.TextSource;

public abstract class AstNode implements IAstNode {
	
	protected TextSource source;
	protected int line;
	protected int offset;
	protected int endoffset;

	public AstNode(TextSource source, int line, int offset, int endoffset) {
		this.source = source;
		this.line = line;
		this.offset = offset;
		this.endoffset = endoffset;
	}

	public int getLine() {
		return this.line;
	}

	public int getOffset() {
		return this.offset;
	}

	public int getEndoffset() {
		return this.endoffset;
	}

	public TextSource getSource() {
		return source;
	}

	public String toString() {
		return source == null ? "" : source.getText(offset, endoffset);
	}

	//public abstract void accept(Visitor v);
}
